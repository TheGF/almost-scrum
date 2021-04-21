package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"time"
)


func pullExchange(connection *Connection, exchange transport.Exchange, since time.Time, r chan Transfer) {
	logrus.Infof("Federation pull from exchange %s", exchange)
	stat := Transfer{Exchange: exchange.Name()}
	start := time.Now()

	list, err := exchange.List(since)
	if err != nil {
		stat.Error = err
		r <- stat
		return
	}
	for _, loc := range list {
		refs, found := connection.locs.LoadOrStore(loc, 1)
		if found {
			connection.locs.Store(loc, refs.(int)+1)
		} else {
			if sz, err := exchange.Pull(loc); err == nil {
				stat.Locs = append(stat.Locs, loc)
				stat.Size += sz
				logrus.Infof("pulled loc %s from %s", loc, exchange)
			} else {
				stat.Issues = append(stat.Issues, err)
				logrus.Warnf("cannot pull %s from %s: %v", loc, exchange, err)
			}
		}
	}
	stat.Elapsed = time.Since(start)
	r <- stat
	return
}

func Pull(project *core.Project, since time.Time) ([]Transfer, error) {
	connection, err := Connect(project)
	if err != nil {
		return nil, err
	}

	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	if since.IsZero() {
		diff := time.Duration(connection.config.Span) * 24 * time.Hour
		since = time.Now().Add(-diff)
	}

	logrus.Infof("Pull project %s with federation node %s starting from %s", project.Config.Public.Name,
		connection.config.UUID, since)

	if err := loadLocs(connection); err != nil {
		return nil, err
	}

	r := make(chan Transfer)
	connected := getConnectedExchanges(connection)
	for _, exchange := range connected {
		go pullExchange(connection, exchange, since, r)
	}

	var stats []Transfer
	for range connected {
		stat := <-r
		stats = append(stats, stat)
		if len(stat.Locs) > 0 {
			connection.throughput[stat.Exchange].Download = stat.Size / int64(len(stat.Locs))
		}
	}


	logrus.Infof("Pull completed for project %s on node %s: %v",
		project.Config.Public.Name, connection.config.UUID, stats)
	return stats, nil
}
