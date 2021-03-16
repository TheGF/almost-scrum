package transport

import (
	"fmt"
	"github.com/ricardopadilha/gowebdav"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type WebDAVConfig struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Timeout  int    `yaml:"timeout"`
}

type WebDAVExchange struct {
	config *WebDAVConfig
	UUID   string
	local  string
	remote string
	conn   *gowebdav.Client
}

func GetWebDAVExchanges(configs ...WebDAVConfig) []Exchange {
	var exchanges []Exchange
	for i := range configs {
		exchanges = append(exchanges, &WebDAVExchange{
			config: &configs[i],
			remote: "",
			local:  "",
			conn:   nil,
		})
	}
	return exchanges
}

func (exchange *WebDAVExchange) ID() string {
	return fmt.Sprintf("webDAV-%s", exchange.config.Name)
}

func (exchange *WebDAVExchange) Connect(remoteRoot, localRoot string) (UpdatesCh, error) {
	if exchange.conn != nil {
		exchange.Disconnect()
	}

	var timeout time.Duration
	if exchange.config.Timeout == 0 {
		timeout = time.Duration(exchange.config.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}

	conn := gowebdav.NewClient(exchange.config.URL, exchange.config.Username, exchange.config.Password)
	conn.SetTimeout(timeout)
	err := conn.Connect()
	if err != nil {
		logrus.Warnf("cannot Connect to WebDAV url %s: %v", exchange.config.URL, err)
		return nil, err
	}

	exchange.remote = remoteRoot
	err = conn.Mkdir(exchange.remote, 0644)
	if err != nil {
		logrus.Warnf("cannot create remote path %s on WebDAV url %s: %v", remoteRoot, exchange.config.URL, err)
		return nil, err
	}

	exchange.local = localRoot
	exchange.remote = remoteRoot
	exchange.conn = conn
	return nil, nil
}

func (exchange *WebDAVExchange) Disconnect() {
	exchange.conn = nil
}

func (exchange *WebDAVExchange) list(folder string, since time.Time) ([]string, error) {
	files, err := exchange.conn.ReadDir(fmt.Sprintf("%s/%s", exchange.remote, folder))
	if err != nil {
		logrus.Warnf("cannot list FTP transport %s: %v",
			exchange.config.URL, err)
	}

	var names []string
	for _, file := range files {
		if file.IsDir() {
			names_, err := exchange.list(path.Join(folder, file.Name()), since)
			if err != nil {
				return nil, err
			}
			names = append(names, names_...)
		} else if file.Size() > 0 && file.ModTime().After(since) {
			names = append(names, path.Join(folder, file.Name()))
		}
	}

	return names, nil
}

func (exchange *WebDAVExchange) List(since time.Time) ([]string, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return nil, os.ErrClosed
	}

	names, err := exchange.list("", since)
	if err == nil {
		logrus.Debugf("list from %s: %#v", exchange, names)
	}
	return names, err
}

func (exchange *WebDAVExchange) Push(loc string) error {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return os.ErrClosed
	}

	c := exchange.conn
	file := filepath.Join(exchange.local, strings.ReplaceAll(loc, "/", string(os.PathSeparator)))
	r, err := os.Open(file)
	if err != nil {
		logrus.Warnf("cannot open loc %s in %s: %v", loc, exchange, err)
		return err
	}
	defer r.Close()

	full := fmt.Sprintf("%s/%s", exchange.remote, loc)
	idx := strings.LastIndex(full, "/")
	dir := full[0:idx]
	info, err := c.Stat(dir)
	if err != nil {
		_ = c.MkdirAll(dir, 0644)
	} else if !info.IsDir() {
		return os.ErrInvalid
	}

	err = c.WriteStream(full, r, 0644)
	if err != nil {
		logrus.Warnf("cannot upload loc %s to %s: %v", loc, exchange, err)
		return err
	}


	logrus.Infof("loc %s uploaded to %s", loc, exchange)
	return nil
}

func (exchange *WebDAVExchange) Pull(loc string) error {
	if exchange.conn == nil {
		logrus.Warn("trying to list without WebDAV connection. Call Connect first")
		return os.ErrClosed
	}

	remote := fmt.Sprintf("%s/%s", exchange.remote, loc)
	r, err := exchange.conn.ReadStream(remote)
	if err != nil {
		return err
	}
	defer r.Close()

	local := filepath.Join(exchange.local, loc)
	_ = os.MkdirAll(filepath.Dir(local), 0755)
	w, err := os.Create(local)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	logrus.Infof("loc %s downloaded from %s to %s", loc, exchange, local)
	return nil
}

func (exchange *WebDAVExchange) String() string {
	return fmt.Sprintf("webDAV %s - %s, connected=%t", exchange.config.Name, exchange.config.URL, exchange.conn != nil)
}
