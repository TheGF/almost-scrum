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
	Name     string `json:"name" yaml:"name"`
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
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

func RemoveWebDAVSecret(configs ...WebDAVConfig) {
	for i := range configs {
		configs[i].Password = ""
	}
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
		logrus.Infof("list from %s: %v", exchange, names)
	}
	return names, err
}

func (exchange *WebDAVExchange) Push(loc string) (int64, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return 0, os.ErrClosed
	}

	c := exchange.conn
	file := filepath.Join(exchange.local, strings.ReplaceAll(loc, "/", string(os.PathSeparator)))
	r, err := os.Open(file)
	if err != nil {
		logrus.Warnf("cannot open loc %s in %s: %v", loc, exchange, err)
		return 0, err
	}
	defer r.Close()

	dest := fmt.Sprintf("%s/%s", exchange.remote, loc)
	idx := strings.LastIndex(dest, "/")
	dir := dest[0:idx]
	stat, err := c.Stat(dir)
	if err != nil {
		_ = c.MkdirAll(dir, 0644)
	} else if !stat.IsDir() {
		return 0, os.ErrInvalid
	}

	err = c.WriteStream(dest, r, 0644)
	if err != nil {
		logrus.Warnf("cannot upload loc %s to %s: %v", loc, exchange, err)
		return 0, err
	}
	stat, _ = os.Stat(file)

	logrus.Infof("loc %s uploaded to %s", loc, exchange)
	return stat.Size(), nil
}

func (exchange *WebDAVExchange) Pull(loc string) (int64, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without WebDAV connection. Call Connect first")
		return 0, os.ErrClosed
	}

	remote := fmt.Sprintf("%s/%s", exchange.remote, loc)
	r, err := exchange.conn.ReadStream(remote)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	local := filepath.Join(exchange.local, loc)
	_ = os.MkdirAll(filepath.Dir(local), 0755)
	w, err := os.Create(local)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	sz, err := io.Copy(w, r)
	if err != nil {
		return 0, err
	}
	logrus.Infof("loc %s downloaded from %s to %s", loc, exchange, local)
	return sz, nil
}

func (exchange *WebDAVExchange) Config(withPrivateKeys bool) interface{} {
	c := *exchange.config
	if !withPrivateKeys {
		c.Password = ""
	}
	return c
}

func (exchange *WebDAVExchange) delete(folder string, pattern string, before time.Time) (int, error) {
	files, err := exchange.conn.ReadDir(path.Join(exchange.remote, folder))
	if err != nil {
		logrus.Warnf("cannot list WebDAV transport %s: %v",
			exchange.config.URL, err)
		return 0, err
	}

	filesCount := 0
	for _, file := range files {
		if file.IsDir() {
			count, err := exchange.delete(path.Join(folder, file.Name()), pattern, before)
			if err != nil {
				return 0, err
			}
			if count == 0 {
				err = exchange.conn.Remove(path.Join(exchange.remote, folder, file.Name()))
				if err != nil {
					logrus.Errorf("cannot remove file %s from %s: %v", file.Name(),
						exchange, err)
				}
			}
		} else {
			match, _ := path.Match(pattern, file.Name())
			if match && file.ModTime().Before(before) {
				if err := exchange.conn.Remove(path.Join(exchange.remote, folder,
					file.Name())); err != nil {
					logrus.Error("cannot remove file %s from %s: %v", file.Name(),
						exchange, err)
				} else {
					logrus.Infof("removed %s from %s", file.Name(), exchange)
				}
			} else {
				filesCount++
			}
		}
	}

	return filesCount, nil
}

func (exchange *WebDAVExchange) Delete(pattern string, before time.Time) error {
	if exchange.conn == nil {
		logrus.Warn("trying to list without WebDAV connection. Call Connect first")
		return os.ErrClosed
	}

	_, err := exchange.delete("", pattern, before)
	return err
}

func (exchange *WebDAVExchange) Name() string {
	return exchange.config.Name
}

func (exchange *WebDAVExchange) String() string {
	return fmt.Sprintf("webDAV %s - %s", exchange.config.Name, exchange.config.URL)
}
