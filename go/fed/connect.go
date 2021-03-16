package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/beevik/ntp"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func getExchanges(config *Config) []transport.Exchange {
	var exchanges []transport.Exchange
	exchanges = append(exchanges, transport.GetFTPExchanges(config.Ftp...)...)
	exchanges = append(exchanges, transport.GetWebDAVExchanges(config.WebDAV...)...)
	return exchanges
}

func loadLocs(signal *Signal) error {
	err := filepath.Walk(signal.local, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			loc, _ := filepath.Rel(signal.local, p)
			signal.locs.LoadOrStore(loc, false)
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("cannot read fed staging folder %s: %v", signal.local, err)
		return err
	}
	return nil
}

func GetSignal(project *core.Project) (*Signal, error) {
	if signal, found := states[project.Config.UUID]; found {
		return signal, nil
	}
	if err := checkTime(); err != nil {
		return nil, err
	}

	localRoot := filepath.Join(project.Path, core.ProjectFedFilesFolder)
	_ = os.MkdirAll(localRoot, 0755)

	config, err := ReadConfig(project)
	if err != nil {
		return nil, err
	}

	signal := &Signal{
		local:     localRoot,
		remote:    project.Config.UUID,
		locs:      sync.Map{},
		exchanges: map[transport.Exchange]bool{},
		config:    config,
	}

	for _, exchange := range getExchanges(config) {
		signal.exchanges[exchange] = false
	}

	if err := loadLocs(signal); err != nil {
		return nil, err
	}

	if connectivity := connect(signal); connectivity == 0 {
		return nil, ErrNoConnections
	}

	go func() {
		signal.reconnect = time.NewTicker(signal.config.ReconnectTime)
		for range signal.reconnect.C {
			connect(signal)
		}
	}()

	logrus.Infof("federation for project %s started successfully", project.Config.UUID)
	states[project.Config.UUID] = signal
	return signal, nil
}

func GetTimeDiff() (int64, error) {
	tm, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return 0, err
	}

	return int64(math.Abs(tm.Sub(time.Now()).Seconds())), nil
}

var timeValid = false

func checkTime() error {
	if timeValid {
		return nil
	}

	for i := 0; i < 10; i++ {
		if diff, err := GetTimeDiff(); err == nil {
			if diff > 60 {
				logrus.Errorf("computer time is not correct. Push is disabled")
				return ErrFedTimeMismatch
			} else {
				timeValid = true
				return nil
			}
		} else {
			logrus.Warnf("cannot connect to time server. Retry %d/10: %v", i, err)
			time.Sleep(2 * time.Second)
		}
	}
	logrus.Errorf("cannot connect to time server")
	return ErrFedNoTimeServer
}

func getConnectedExchanges(signal* Signal) []transport.Exchange{
	var exchanges []transport.Exchange
	for exchange, connected := range signal.exchanges {
		if connected {
			exchanges = append(exchanges, exchange)
		}
	}
	return exchanges
}

func asyncUpdates(signal *Signal, exchange transport.Exchange, ch transport.UpdatesCh) {
	for {
		locs, open := <-ch
		if !open {
			break
		}
		signal.inUse.Add(1)
		for _, loc := range locs {
			if _, loaded := signal.locs.LoadOrStore(loc, true); !loaded {
				if err := exchange.Pull(loc); err != nil {
					signal.locs.Delete(loc)
				}
			}
		}
		signal.inUse.Done()
	}
}

func connect(signal *Signal) int {
	c := make(chan bool)

	var disconnected []transport.Exchange
	for exchange, connected := range signal.exchanges {
		if !connected {
			disconnected = append(disconnected, exchange)
		}
	}

	for _, exchange := range disconnected {
		go func(exchange transport.Exchange) {
			ch, err := exchange.Connect(signal.remote, signal.local)
			if err == nil {
				logrus.Infof("connected to transport %s", exchange)
				c <- true
				signal.exchanges[exchange] = true
				if ch != nil {
					go asyncUpdates(signal, exchange, ch)
				}
			} else {
				c <- false
				logrus.Infof("failed connect to transport %s: %v", exchange, err)
			}
		}(exchange)
	}

	newConnected := 0
	for range disconnected {
		if <-c {
			newConnected++
		}
	}

	return len(signal.exchanges)-len(disconnected)+newConnected
}

func Disconnect(project *core.Project) {
	if state, found := states[project.Config.UUID]; found {
		state.reconnect.Stop()
		state.inUse.Wait()
		for hub, connected := range state.exchanges {
			if connected {
				hub.Disconnect()
				state.exchanges[hub] = false
			}
		}
	}
}
