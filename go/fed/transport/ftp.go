package transport

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type FTPConfig struct {
	Name     string `json:"name" yaml:"name"`
	URL      string `json:"url" yaml:"url"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Timeout  int    `json:"timeout" yaml:"timeout"`
}

type FTPExchange struct {
	config *FTPConfig
	remote string
	local  string
	conn   *ftp.ServerConn
}

func GetFTPExchanges(configs ...FTPConfig) []Exchange {
	var exchanges []Exchange
	for i := range configs {
		exchanges = append(exchanges, &FTPExchange{
			config: &configs[i],
			remote: "",
			local:  "",
			conn:   nil,
		})
	}
	return exchanges
}

func RemoveFTPSecret(configs ...FTPConfig) {
	for i := range configs {
		configs[i].Password = ""
	}
}

func (exchange *FTPExchange) ID() string {
	return fmt.Sprintf("ftp-%s", exchange.config.Name)
}

func (exchange *FTPExchange) Connect(remoteRoot, localRoot string) (UpdatesCh, error) {
	if exchange.conn != nil {
		exchange.Disconnect()
	}

	url_, err := url.Parse(exchange.config.URL)
	if err != nil {
		logrus.Warnf("FTP url %s is invalid: %v", exchange.config.URL, err)
	}

	var timeout time.Duration
	if exchange.config.Timeout == 0 {
		timeout = time.Duration(exchange.config.Timeout) * time.Second
	} else {
		timeout = 30 * time.Second
	}

	conn, err := ftp.Dial(url_.Host, ftp.DialWithTimeout(timeout))
	if err != nil {
		logrus.Warnf("cannot Connect to FTP transport %s: %v", exchange.config.URL, err)
		return nil, err
	}
	err = conn.Login(exchange.config.Username, exchange.config.Password)
	if err != nil {
		logrus.Warnf("cannot authenticate to FTP transport %s: %v", exchange.config.URL, err)
		return nil, err
	}

	p := fmt.Sprintf("%s/%s", url_.Path, remoteRoot)
	_ = conn.MakeDir(p)
	err = conn.ChangeDir(p)
	if err != nil {
		logrus.Warnf("cannot change to local %s on FTP host %s: %v", remoteRoot, url_.Host, err)
		return nil, err
	}

	exchange.local = localRoot
	exchange.remote = remoteRoot
	exchange.conn = conn
	return nil, nil
}

func (exchange *FTPExchange) Disconnect() {
	if exchange.conn != nil {
		_ = exchange.conn.Quit()
		exchange.conn = nil
	}
}

func (exchange *FTPExchange) list(folder string, since time.Time) ([]string, error) {
	entries, err := exchange.conn.List(folder)
	if err != nil {
		logrus.Warnf("cannot list FTP transport %s: %v",
			exchange.config.URL, err)
	}

	var names []string
	for _, entry := range entries {

		switch entry.Type {
		case ftp.EntryTypeFile:
			if entry.Size > 0 && entry.Time.After(since) {
				names = append(names, path.Join(folder,entry.Name))
			}
		case ftp.EntryTypeFolder:
			names_, err := exchange.list(path.Join(folder, entry.Name), since)
			if err != nil {
				return nil, err
			}
			names = append(names, names_...)
		}
	}

	return names, nil
}

func (exchange *FTPExchange) List(since time.Time) ([]string, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return nil, os.ErrClosed
	}

	names, err := exchange.list("", since)
	if err == nil {
		logrus.Debugf("list from %s: %v", exchange, names)
	}
	return names, err
}

func (exchange *FTPExchange) Push(loc string) (int64, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return 0, os.ErrClosed
	}

	c := exchange.conn
	file := filepath.Join(exchange.local, strings.ReplaceAll(loc, "/", string(os.PathSeparator)))
	stat, _ := os.Stat(file)
	r, err := os.Open(file)
	if err != nil {
		logrus.Warnf("cannot open loc %s in %s: %v", loc, exchange, err)
		return 0, err
	}
	defer r.Close()

	//	l := filepath.Base(loc)

	curr, _ := c.CurrentDir()
	parts := strings.Split(loc, "/")
	for i := 0; i < len(parts)-1; i++ {
		_ = c.MakeDir(parts[i])
		_ = c.ChangeDir(parts[i])
	}

	err = exchange.conn.Stor(parts[len(parts)-1], r)
	_ = c.ChangeDir(curr)
	if err != nil {
		logrus.Warnf("cannot upload loc %s to %s: %v", loc, exchange, err)
		return 0, err
	}

	logrus.Infof("loc %s uploaded to %s", loc, exchange)
	return stat.Size(), nil
}

func (exchange *FTPExchange) Pull(file string)  (int64, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return 0, os.ErrClosed
	}

	r, err := exchange.conn.Retr(file)
	if err != nil {
		return 0, err
	}
	defer r.Close()

	dest := filepath.Join(exchange.local, file)
	_ = os.MkdirAll(filepath.Dir(dest), 0755)
	w, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer w.Close()

	sz, err := io.Copy(w, r)
	if err != nil {
		return 0, err
	}
	logrus.Infof("file %s downloaded from %s to %s", file, exchange, dest)
	return sz, nil
}


func (exchange *FTPExchange) delete(folder string, pattern string, before time.Time) (int, error) {
	entries, err := exchange.conn.List(folder)
	if err != nil {
		logrus.Warnf("cannot list FTP transport %s: %v",
			exchange.config.URL, err)
	}

	filesCount := 0
	for _, entry := range entries {
		p := path.Join(folder, entry.Name)
		switch entry.Type {
		case ftp.EntryTypeFile:
			match, _ := path.Match(pattern, entry.Name)
			if match && entry.Time.Before(before) {
				if err := exchange.conn.Delete(p); err != nil  {
					logrus.Errorf("cannot remove %s from %s: %v", entry.Name, exchange, err)
				} else {
					logrus.Infof("removed %s from %s", entry.Name, exchange)
				}
			} else {
				filesCount++
			}
		case ftp.EntryTypeFolder:
			count, err := exchange.delete(p, pattern, before)
			if err != nil {
				return 0, err
			}
			if count == 0 {
				exchange.conn.Delete(p)
			}
		}
	}

	return filesCount, nil
}


func (exchange *FTPExchange) Delete(pattern string, before time.Time) error {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return os.ErrClosed
	}

	_, err := exchange.delete("", pattern, before)
	return err
}


func (exchange *FTPExchange) Name() string {
	return exchange.config.Name
}

func (exchange *FTPExchange) String() string {
	return fmt.Sprintf("ftp %s - %s", exchange.config.Name, exchange.config.URL)
}
