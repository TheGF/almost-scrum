package fed

import (
	"almost-scrum/core"
	"almost-scrum/fed/transport"
	"github.com/sirupsen/logrus"
	"time"
)

func pushExchange(connection *Connection, exchange transport.Exchange, r chan Transfer) {
	logrus.Debugf("Federation push to exchange %s", exchange)
	stat := Transfer{Exchange: exchange.Name()}
	start := time.Now()

	remotes, err := exchange.List(time.Time{})
	if err != nil {
		stat.Error = err
		r <- stat
		return
	}

	connection.locs.Range(func(key, refs interface{}) bool {
		loc := key.(string)
		if _, found := core.FindStringInSlice(remotes, loc); !found {

			logrus.Debugf("location %s is not present in %s", loc, exchange)
			if sz, err := exchange.Push(loc); err == nil {
				logrus.Infof("location %s uploaded to %s", loc, exchange)
				connection.locs.Store(loc, refs.(int)+1)
				stat.Locs = append(stat.Locs, loc)
				stat.Size += sz
			} else {
				stat.Issues = append(stat.Issues, err)
				logrus.Warnf("cannot push %s to %s: %v", loc, exchange, err)
			}
		} else {
			logrus.Debugf("location %s is present in %s", loc, exchange)
		}
		return true
	})
	stat.Elapsed = time.Since(start)
	r <- stat
	return
}

func Push(project *core.Project) ([]Transfer, error) {
	connection, err := Connect(project)
	if err != nil {
		return nil, err
	}

	connection.mutex.Lock()
	defer connection.mutex.Unlock()

	logrus.Debugf("Push project %s with federation node %s", project.Config.Public.Name,
		connection.config.UUID)

	if err := loadLocs(connection); err != nil {
		return nil, err
	}

	r := make(chan Transfer)
	connected := getConnectedExchanges(connection)
	for _, exchange := range connected {
		go pushExchange(connection, exchange, r)
	}

	var stats []Transfer
	for range connected {
		stat := <-r
		stats = append(stats, stat)
		if len(stat.Locs) > 0 {
			connection.throughput[stat.Exchange].Upload = stat.Size / int64(len(stat.Locs))
		}
	}

	logrus.Debugf("Push completed for project %s on node %s: %v",
		project.Config.Public.Name, connection.config.UUID, stats)
	return stats, nil
}
