package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Status struct {
	Exchanges map[string]bool `json:"exchanges"`
	Locs      *sync.Map        `json:"locs"`
	Stat      map[string]Stat `json:"stat"`
}

func GetStatus(project *core.Project) Status {
	signal, err := Connect(project)
	if err != nil {
		return Status{}
	}

	s := Status{
		Exchanges: map[string]bool{},
		Locs:      &signal.locs,
		Stat:      map[string]Stat{},
	}
	for exchange, connected := range signal.exchanges {
		s.Exchanges[exchange.Name()] = connected
	}
	for exchange, stat := range signal.stat {
		s.Stat[exchange.Name()] = stat
	}

	return s
}

func pushExchange(signal *Connection, exchange transport.Exchange, remote map[string]bool, stat *Stat) (errs []error) {
	start := time.Now()
	signal.locs.Range(func(loc, refs interface{}) bool {
		if _, found := remote[loc.(string)]; !found {
			if sz, err := exchange.Push(loc.(string)); err == nil {
				logrus.Infof("push %s to %s", loc.(string), exchange)
				signal.locs.Store(loc, refs.(int)+1)
				stat.Upload += sz
				stat.Push++
			} else {
				errs = append(errs, err)
				logrus.Warnf("cannot push %s to %s: %v", loc, exchange, err)
			}
		} else {
			remote[loc.(string)] = true
		}
		return true
	})
	stat.Upload = stat.Upload * int64(time.Second) / int64(time.Since(start))
	return
}

func pullExchange(signal *Connection, exchange transport.Exchange, remote map[string]bool, stat *Stat) (errs []error) {
	start := time.Now()
	for loc, present := range remote {
		if !present {
			if refs, found := signal.locs.LoadOrStore(loc, 1); !found {
				if sz, err := exchange.Pull(loc); err == nil {
					stat.Download += sz
					stat.Pull++
				} else {
					errs = append(errs, err)
					signal.locs.Delete(loc)
					logrus.Warnf("cannot pull %s from %s: %v", loc, exchange, err)
				}
			} else {
				logrus.Infof("pull %s from %s", loc, exchange)
				signal.locs.Store(loc, refs.(int)+1)
			}
		}
	}
	stat.Download = stat.Download * int64(time.Second) / int64(time.Since(start))
	return
}

func syncExchange(connection *Connection, exchange transport.Exchange, since time.Time, r chan bool) {
	logrus.Infof("Federation sync from exchange %s", exchange)

	list, err := exchange.List(since)
	if err != nil {
		r <- false
		return
	}

	remote := map[string]bool{}
	for _, loc := range list {
		remote[loc] = false
	}

	stat := Stat{}
	ok := len(pullExchange(connection, exchange, remote, &stat)) == 0
	ok = ok && len(pushExchange(connection, exchange, remote, &stat)) == 0
	connection.stat[exchange] = stat

	diff := 3 * time.Duration(connection.config.Span) * 24 * time.Hour
	before := time.Now().Add(-diff)
	_ = exchange.Delete("dat", before)

	r <- ok
}

func Sync(project *core.Project, since time.Time) (failedExchanges int, err error) {
	connection, err := Connect(project)
	if err != nil {
		return 0, err
	}

	if since.IsZero() {
		diff := time.Duration(connection.config.Span) * 24 * time.Hour
		since = time.Now().Add(-diff)
	}

	now := time.Now()
	if diff := time.Now().Sub(connection.config.LastSync).Seconds(); diff < 30 {
		logrus.Infof("last sync only %f seconds ago; try later", diff)
		return 0, nil
	}
	connection.config.LastSync = now

	logrus.Infof("Synchronize project %s with federation node %s starting from %s", project.Config.Public.Name,
		connection.config.UUID, since)

	if err := loadLocs(connection); err != nil {
		return 0, err
	}

	r := make(chan bool)
	connected := getConnectedExchanges(connection)
	for _, exchange := range connected {
		go syncExchange(connection, exchange, since, r)
	}

	for range connected {
		if !<-r {
			failedExchanges++
		}
	}

//	_ = WriteConfig(project, connection.config)

	logrus.Infof("Synchronization completed for project %s on node %s with %d failed exchanges",
		project.Config.Public.Name, connection.config.UUID,	failedExchanges)
	return
}
