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
	exchanges = append(exchanges, transport.GetS3Exchanges(config.S3...)...)
	exchanges = append(exchanges, transport.GetWebDAVExchanges(config.WebDAV...)...)
	exchanges = append(exchanges, transport.GetFTPExchanges(config.Ftp...)...)
	return exchanges
}

func loadLocs(signal *Connection) error {
	locs := sync.Map{}
	count := 0
	err := filepath.Walk(signal.local, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			loc, _ := filepath.Rel(signal.local, p)
			locs.Store(loc, 0)
			count++
		}
		return nil
	})
	if err != nil {
		logrus.Errorf("cannot read fed staging folder %s: %v", signal.local, err)
		return err
	}

	signal.locs = locs
	logrus.Infof("loaded %d files in matching map", count)
	return nil
}

func Connect(project *core.Project) (*Connection, error) {
	if signal, found := states[project.Config.UUID]; found {
		return signal, nil
	}
	if err := checkTime(); err != nil {
		return nil, err
	}

	localRoot := filepath.Join(project.Path, core.ProjectFedFilesFolder)
	_ = os.MkdirAll(localRoot, 0755)

	config, err := ReadConfig(project, false)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		local:     localRoot,
		remote:    config.UUID,
		locs:      sync.Map{},
		exchanges: map[transport.Exchange]bool{},
		config:    config,
		stat:      map[transport.Exchange]Stat{},
	}

	for _, exchange := range getExchanges(config) {
		connection.exchanges[exchange] = false
	}

	if err := loadLocs(connection); err != nil {
		return nil, err
	}

	open(connection)

	go func() {
		connection.reconnect = time.NewTicker(connection.config.ReconnectTime)
		for range connection.reconnect.C {
			open(connection)
		}
	}()

	logrus.Infof("federation for project %s started successfully", project.Config.UUID)
	states[project.Config.UUID] = connection
	return connection, nil
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
			logrus.Warnf("cannot open to time server. Retry %d/10: %v", i, err)
			time.Sleep(2 * time.Second)
		}
	}
	logrus.Errorf("cannot open to time server")
	return ErrFedNoTimeServer
}

func getConnectedExchanges(signal *Connection) []transport.Exchange {
	var exchanges []transport.Exchange
	for exchange, connected := range signal.exchanges {
		if connected {
			exchanges = append(exchanges, exchange)
		}
	}
	return exchanges
}

func asyncUpdates(signal *Connection, exchange transport.Exchange, ch transport.UpdatesCh) {
	for {
		locs, open := <-ch
		if !open {
			break
		}
		signal.inUse.Add(1)
		for _, loc := range locs {
			if _, loaded := signal.locs.LoadOrStore(loc, true); !loaded {
				if _, err := exchange.Pull(loc); err != nil {
					signal.locs.Delete(loc)
				}
			}
		}
		signal.inUse.Done()
	}
}

func open(connection *Connection) int {
	c := make(chan bool)

	var disconnected []transport.Exchange
	for exchange, connected := range connection.exchanges {
		if !connected {
			disconnected = append(disconnected, exchange)
		}
	}

	for _, exchange := range disconnected {
		go func(exchange transport.Exchange) {
			ch, err := exchange.Connect(connection.remote, connection.local)
			if err == nil {
				logrus.Infof("connected to transport %s", exchange)
				c <- true
				connection.exchanges[exchange] = true
				if ch != nil {
					go asyncUpdates(connection, exchange, ch)
				}
			} else {
				c <- false
				logrus.Infof("failed open to transport %s: %v", exchange, err)
			}
		}(exchange)
	}

	newConnected := 0
	for range disconnected {
		if <-c {
			newConnected++
		}
	}

	return len(connection.exchanges) - len(disconnected) + newConnected
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
