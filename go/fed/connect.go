package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"errors"
	"github.com/beevik/ntp"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"math"
	"net"
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
	exchanges = append(exchanges, transport.GetUSBExchanges(config.USB...)...)
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
	return nil
}

func Connect(project *core.Project) (*Connection, error) {

	if connection, found := connections.Get(project.Config.UUID); found {
		return connection.(*Connection), nil
	}

	localRoot := filepath.Join(project.Path, core.ProjectFedFilesFolder)
	_ = os.MkdirAll(localRoot, 0755)

	config, err := ReadConfig(project, false)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		local:      localRoot,
		remote:     config.UUID,
		locs:       sync.Map{},
		exchanges:  map[transport.Exchange]bool{},
		config:     config,
		throughput: map[string]*Throughput{},
		exports:    map[string]time.Time{},
		checkTime:  checkTime(),
	}

	for _, exchange := range getExchanges(config) {
		connection.exchanges[exchange] = false
		connection.throughput[exchange.Name()] = &Throughput{}
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
	connections.Set(project.Config.UUID, connection, cache.DefaultExpiration)
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
		var networkError *net.UnknownNetworkError
		if diff, err := GetTimeDiff(); err == nil {
			if diff > 60 {
				logrus.Errorf("computer time is not correct. Push is disabled")
				return ErrFedTimeMismatch
			} else {
				timeValid = true
				return nil
			}
		} else if errors.As(err, &networkError) {
			logrus.Warnf("No connection to ModTime server: %v", err)
			return err
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

func asyncUpdates(connection *Connection, exchange transport.Exchange, ch transport.UpdatesCh) {
	for {
		locs, open := <-ch
		if !open {
			break
		}
		connection.mutex.Lock()
		defer connection.mutex.Unlock()

		for _, loc := range locs {
			if _, loaded := connection.locs.LoadOrStore(loc, true); !loaded {
				if _, err := exchange.Pull(loc); err != nil {
					connection.locs.Delete(loc)
				}
			}
		}
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

	if c, found := connections.Get(project.Config.UUID); found {
		connection := c.(*Connection)
		connection.reconnect.Stop()

		connection.mutex.Lock()
		defer connection.mutex.Unlock()

		for hub, connected := range connection.exchanges {
			if connected {
				hub.Disconnect()
				connection.exchanges[hub] = false
			}
		}
	}
}
