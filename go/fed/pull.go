package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"time"
)


func pullFromExchange(signal *Signal, exchange transport.Exchange, since time.Time, r chan error)  {
	go func() {
		logrus.Infof("start pull from transport %s for locs newer than %s", exchange, since)
		locs, err := exchange.List(since)
		if err != nil {
			logrus.Errorf("error during list from transport %s: %v", exchange, err)
			r <- err
		}
		if len(locs) == 0{
			logrus.Infof("no locs from transport %s", exchange)
			r <- err
		}

		var imported []string
		for _, loc := range locs {
			if _, loaded := signal.locs.LoadOrStore(loc, true); !loaded {
				if _, err := exchange.Pull(loc); err == nil {
					logrus.Debugf("import %s from transport %s", loc, exchange)
					imported = append(imported, loc)
				} else {
					signal.locs.Delete(loc)
				}
			}
		}
		logrus.Infof("end import from transport %s, imported locs: %#v", exchange, imported)
		r <- nil
	}()
}


func Pull(project *core.Project, since time.Time) (int, error) {
	signal, err := GetSignal(project)
	if err != nil {
		return 0, err
	}
	signal.inUse.Add(1)
	defer signal.inUse.Done()

	r := make(chan error, 8)
	defer close(r)
	connected := getConnectedExchanges(signal)
	for _, exchange := range connected {
		go pullFromExchange(signal, exchange, since, r)
	}

	successful := 0
	for range connected{
		if <- r == nil {
			successful++
		}
	}

	return successful, nil
}
