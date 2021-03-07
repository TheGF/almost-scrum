package federation

import (
	"almost-scrum/core"
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type FTPConfig struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type FTPHub struct {
	Config *FTPConfig
	UUID string
	Conn   *ftp.ServerConn
}

func getFTPHubs(project* core.Project, config *Config) []Hub {
	var hubs []Hub
	for _, c := range config.Ftp {
		hubs = append(hubs, &FTPHub{
			Config: &c,
			UUID: project.Config.UUID,
			Conn:   nil,
		})
	}
	return hubs
}

func (hub *FTPHub) ID() string {
	return fmt.Sprintf("ftp-%s",hub.Config.Name)
}

func (hub *FTPHub) Connect() error {
	if hub.Conn != nil {
		hub.Disconnect()
	}

	url_, err := url.Parse(hub.Config.URL)
	if err != nil {
		logrus.Warnf("FTP url %s is invalid: %v", hub.Config.URL, err)
	}

	conn, err := ftp.Dial(url_.Host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		logrus.Warnf("cannot connect to FTP hub %s: %v", hub.Config.URL, err)
		return err
	}
	err = conn.Login(hub.Config.Username, hub.Config.Password)
	if err != nil {
		logrus.Warnf("cannot authenticate to FTP hub %s: %v", hub.Config.URL, err)
		return err
	}

	path := fmt.Sprintf("%s/%s", url_.Path, hub.UUID)
	_ = conn.MakeDir(path)
	err = conn.ChangeDir(path)
	if err != nil {
		logrus.Warnf("cannot change to path %s on FTP host %s: %v", path, url_.Host, err)
		return err
	}

	hub.Conn = conn
	return nil
}

func (hub *FTPHub) Disconnect() {
	if hub.Conn != nil {
		_ = hub.Conn.Quit()
		hub.Conn = nil
	}
}

func (hub *FTPHub) List(time time.Time) ([]string, error) {
	if hub.Conn == nil {
		logrus.Warn("trying to list without FTP connection. Call connect first")
		return nil, os.ErrClosed
	}

	entries, err := hub.Conn.List("")
	if err != nil {
		logrus.Warnf("cannot list FTP hub %s: %v",
			hub.Config.URL, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.Type == ftp.EntryTypeFile && entry.Time.After(time){
			names = append(names, entry.Name)
		}
	}

	return names, nil
}
func (hub *FTPHub) Push(file string) error {
	if hub.Conn == nil {
		logrus.Warn("trying to list without FTP connection. Call connect first")
		return os.ErrClosed
	}

	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	name := filepath.Base(file)
	err = hub.Conn.Stor(name, r)
	if err != nil {
		return err
	}

	return nil
}


func (hub *FTPHub) Pull(name string, path string) error {
	if hub.Conn == nil {
		logrus.Warn("trying to list without FTP connection. Call connect first")
		return os.ErrClosed
	}

	r, err := hub.Conn.Retr(name)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(filepath.Join(path, name))
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}


func (hub *FTPHub) String() string {
	return fmt.Sprintf("ftp %s - %s, connected=%t", hub.Config.Name, hub.Config.URL, hub.Conn != nil)
}
