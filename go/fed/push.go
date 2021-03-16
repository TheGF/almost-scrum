package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"time"
)

func pushToExchange(signal *Signal, exchange transport.Exchange, r chan error) {
	span := time.Duration(signal.config.Span) * time.Hour * 24
	since := time.Now().Add(-span)
	files, err := exchange.List(since)
	if err != nil {
		logrus.Errorf("error during list from transport %s: %v", exchange, err)
		r <- err
	}

	signal.locs.Range(func(loc, _ interface{}) bool {
		_, found := core.FindStringInSlice(files, loc.(string))
		if !found {
			err := exchange.Push(loc.(string))
			r <- err
			if err != nil {
				logrus.Warnf("cannot push file %s to %s: %v", loc, exchange, err)
				return false
			}
		}
		return true
	})
}

func Push(project *core.Project) (int, error) {
	signal, _ := GetSignal(project)
	signal.inUse.Add(1)
	defer signal.inUse.Done()

	if err := loadLocs(signal); err != nil {
		return 0, err
	}

	r := make(chan error)
	defer close(r)
	connected := getConnectedExchanges(signal)
	for _, exchange := range connected {
		go pushToExchange(signal, exchange, r)
	}

	successful := 0
	for range connected{
		if <- r == nil {
			successful++
		}
	}

	return successful, nil
}
