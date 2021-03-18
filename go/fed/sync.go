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
	Locs      sync.Map        `json:"locs"`
	Stat      map[string]Stat `json:"stat"`
}

func GetStatus(project *core.Project) Status {
	signal, err := GetSignal(project)
	if err != nil {
		return Status{}
	}

	s := Status{
		Exchanges: map[string]bool{},
		Locs:      signal.locs,
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

func syncExchange(signal *Signal, exchange transport.Exchange, since time.Time, r chan error) {
	logrus.Infof("Federation sync from exchange %s", exchange)

	list, err := exchange.List(since)
	if err != nil {
		r <- err
		return
	}

	remote := map[string]bool{}
	for _, loc := range list {
		remote[loc] = false
	}

	stat := Stat{}

	var lastErr error
	// Push files
	start := time.Now()
	signal.locs.Range(func(loc, refs interface{}) bool {
		if _, found := remote[loc.(string)]; !found {
			if sz, err := exchange.Push(loc.(string)); err == nil {
				logrus.Infof("push %s to %s", loc.(string), exchange)
				signal.locs.Store(loc, refs.(int)+1)
				stat.Upload += sz
				stat.Push++
			} else {
				lastErr = err
				logrus.Warnf("cannot push %s to %s: %v", loc, exchange, err)
			}
		} else {
			remote[loc.(string)] = true
		}
		return true
	})
	stat.Upload = stat.Upload * int64(time.Second) / int64(time.Since(start))

	// Pull files
	start = time.Now()
	for loc, present := range remote {
		if !present {
			if refs, found := signal.locs.LoadOrStore(loc, 1); !found {
				if sz, err := exchange.Pull(loc); err == nil {
					stat.Download += sz
					stat.Pull++
				} else {
					lastErr = err
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

	signal.stat[exchange] = stat
	r <- lastErr
}

func Sync(project *core.Project, since time.Time) (failedExchanges int, err error) {
	signal, err := GetSignal(project)
	if err != nil {
		return 0, err
	}


	if since.IsZero() {
		diff := time.Duration(signal.config.Span) * 24 * time.Hour
		since = time.Now().Add(-diff)
	}

	logrus.Infof("Synchronize project %s with federation starting from %s", project.Config.UUID, since)

	if err := loadLocs(signal); err != nil {
		return 0, err
	}

	r := make(chan error)
	connected := getConnectedExchanges(signal)
	for _, exchange := range connected {
		go syncExchange(signal, exchange, since, r)
	}

	for range connected {
		if <-r != nil {
			failedExchanges++
		}
	}
	logrus.Infof("Synchronization completed for project %s with %d failed exchanges", project.Config.UUID,
		failedExchanges)
	return
}
