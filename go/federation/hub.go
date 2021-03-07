package federation

import (
	"almost-scrum/core"
	"errors"
	"github.com/beevik/ntp"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

var (
	ErrFedNoTimeServer = errors.New("no time server")
	ErrFedTimeMismatch = errors.New("time mismatch")
	ErrFedCorrupted    = errors.New("corrupted file")
)

const FedHeaderFile = "ash-header.json"

type Hub interface {
	ID() string
	Connect() error
	Disconnect()

	List(time time.Time) ([]string, error)
	Push(file string) error
	Pull(file string, path string) error

	String() string
}

type Header struct {
	ProjectID string
	PeerID    string
	Peer      string
	Time      time.Time
	User      string
}


func getHubs(project *core.Project) ([]Hub, *Config, error) {
	var hubs []Hub
	config, err := ReadConfig(project)
	if err != nil {
		logrus.Errorf("cannot read federation config: %v", err)
		return nil, nil, err
	}

	hubs = append(hubs, getFTPHubs(project, config)...)
	return hubs, config, nil
}

func connectToHubs(project *core.Project) ([]Hub, *Config, error) {
	hubs, config, err := getHubs(project)
	if err != nil {
		return nil, nil, err
	}

	var connected []Hub
	for _, hub := range hubs {
		if err := hub.Connect(); err == nil {
			connected = append(connected, hub)
			logrus.Infof("connected to hub %s", hub)
		} else {
			logrus.Warnf("connection to hub %s not available: %v", hub, err)
		}
	}
	return connected, config, nil
}



func GetTimeDiff() (int64, error) {
	tm, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return 0, err
	}

	return int64(math.Abs(tm.Sub(time.Now()).Seconds())), nil
}

func CheckTime() error {
	if diff, err := GetTimeDiff(); err != nil {
		logrus.Errorf("cannot connect to time server: %v", err)
		return ErrFedNoTimeServer
	} else if diff > 60 {
		logrus.Errorf("computer time is not correct. Export is disabled")
		return ErrFedTimeMismatch
	} else {
		return nil
	}
}