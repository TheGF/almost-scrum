package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"time"
)

type Status struct {
	exchanges map[string]bool
	outgoing  []string
}

func GetStatus(project *core.Project) Status {
	signal, err := GetSignal(project)
	if err != nil {
		return Status{}
	}

	s := Status{
		exchanges:    map[string]bool{},
		outgoing: nil,
	}

	for exchange, connected := range signal.exchanges {
		s.exchanges[exchange.String()] = connected
	}
	signal.locs.Range(func(loc, uploaded interface{}) bool {
		if !uploaded.(bool) {
			s.outgoing = append(s.outgoing, loc.(string))
		}
		return true
	})

	return s
}

func syncExchange(signal *Signal, exchange transport.Exchange, since time.Time, r chan error) {
	list, err := exchange.List(since)
	if err != nil {
		r <- err
		return
	}

	remote := map[string]bool{}
	for _, loc := range list {
		remote[loc] = false
	}

	var lastErr error
	// Push files
	signal.locs.Range(func(loc, sync interface{}) bool {
		logrus.Infof("should push %s to %s", loc.(string), exchange)
		if _, found := remote[loc.(string)]; !found {
			logrus.Infof("try to push %s to %s", loc.(string), exchange)
			if err := exchange.Push(loc.(string)); err != nil {
				lastErr = err
				logrus.Warnf("cannot push %s to %s: %v", loc, exchange, err)
			}
		} else {
			remote[loc.(string)] = true
		}
		return true
	})

	// Pull files
	for loc, present := range remote {
		if !present {
			if _, found := signal.locs.LoadOrStore(loc, true); !found {
				logrus.Infof("try to pull %s from %s", loc, exchange)
				if err := exchange.Pull(loc); err != nil {
					lastErr = err
					signal.locs.Delete(loc)
					logrus.Warnf("cannot pull %s from %s: %v", loc, exchange, err)
				}
			}
		}
	}
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

	loadLocs(signal)

	r := make(chan error)
	connected := getConnectedExchanges(signal)
	for _, exchange := range connected {
		go syncExchange(signal, exchange, since, r)
	}

	for range connected {
		if <- r != nil {
			failedExchanges++
		}
	}
	return
}
