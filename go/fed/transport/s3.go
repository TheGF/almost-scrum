package transport

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)


type S3Config struct {
	Name      string `json:"name" yaml:"name"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	Bucket    string `json:"bucket" yaml:"bucket"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	Secret    string `json:"secret" yaml:"secret"`
	UseSSL    bool   `json:"useSSL" yaml:"useSSL"`
	Location  string `json:"location" yaml:"location"`
}

type S3Exchange struct {
	config *S3Config
	remote string
	local  string
	conn   *minio.Client
}

func GetS3Exchanges(configs ...S3Config) []Exchange {
	var exchanges []Exchange
	for i := range configs {
		exchanges = append(exchanges, &S3Exchange{
			config: &configs[i],
			remote: "",
			local:  "",
			conn:   nil,
		})
	}
	return exchanges
}

func RemoveS3Secret(configs ...S3Config) {
	for i := range configs {
		configs[i].Secret = ""
	}
}


func (exchange *S3Exchange) ID() string {
	return fmt.Sprintf("ftp-%s", exchange.config.Name)
}

func (exchange *S3Exchange) Connect(remoteRoot, localRoot string) (UpdatesCh, error) {
	if exchange.conn != nil {
		exchange.Disconnect()
	}

	ctx := context.Background()
	defer ctx.Done()

	c := exchange.config
	conn, err := minio.New(c.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.AccessKey, c.Secret, ""),
		Secure: c.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	location := c.Location
	if location == "" {
		location = "us-east-1"
	}

	exists, err := conn.BucketExists(ctx, c.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = conn.MakeBucket(ctx, c.Bucket, minio.MakeBucketOptions{Region: location})
		if err != nil {
			logrus.Errorf("cannot create bucket %s on S3 exchange %s (%s)", c.Bucket,
				c.Endpoint, location)
			return nil, err
		}
	}

	exchange.conn = conn
	exchange.local = localRoot
	exchange.remote = remoteRoot
	return nil, nil
}

func (exchange *S3Exchange) Disconnect() {
	if exchange.conn != nil {
		exchange.conn = nil
	}
}

func (exchange *S3Exchange) List(since time.Time) ([]string, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return nil, os.ErrClosed
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entries := exchange.conn.ListObjects(ctx, exchange.config.Bucket, minio.ListObjectsOptions{
		Prefix:    exchange.remote,
		Recursive: true,
	})

	cut := len(exchange.remote)+1
	var names []string
	for entry := range entries {
		name := entry.Key
		if !strings.HasSuffix(name, "/") && entry.Size > 0 && entry.LastModified.After(since) {
			names = append(names, name[cut:])
		}
	}

	logrus.Infof("list from %s: %v", exchange, names)
	return names, nil
}

func (exchange *S3Exchange) Push(loc string) (int64, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return 0, os.ErrClosed
	}

	file := filepath.Join(exchange.local, strings.ReplaceAll(loc, "/", string(os.PathSeparator)))

	mime, _ := mimetype.DetectFile(file)
	name := path.Join(exchange.remote, loc)
	uploadInfo, err := exchange.conn.FPutObject(context.Background(), exchange.config.Bucket, name,
		file, minio.PutObjectOptions{
			ContentType: mime.String(),
		})
	if err != nil {
		logrus.Warnf("cannot upload loc %s to %s: %v", loc, exchange, err)
		return 0, err
	}

	logrus.Infof("loc %s uploaded to %s: %v", loc, exchange, uploadInfo)
	return uploadInfo.Size, nil
}

func (exchange *S3Exchange) Pull(loc string) (int64, error) {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return 0, os.ErrClosed
	}

	dest := filepath.Join(exchange.local, loc)
	_ = os.MkdirAll(filepath.Dir(dest), 0755)

	name := path.Join(exchange.remote, loc)
	err := exchange.conn.FGetObject(context.Background(), exchange.config.Bucket, name, dest,
		minio.GetObjectOptions{})
	if err != nil {
		logrus.Warnf("cannot download loc %s from %s: %v", loc, exchange, err)
		return 0, err
	}
	stat, _ := os.Stat(dest)

	logrus.Infof("loc %s downloaded from %s to %s", loc, exchange, dest)
	return stat.Size(), nil
}

func (exchange *S3Exchange) Config(withPrivateKeys bool) interface{} {
	c := *exchange.config
	if !withPrivateKeys {
		c.Secret = ""
	}
	return c
}

func (exchange *S3Exchange) Delete(pattern string, before time.Time) error {
	if exchange.conn == nil {
		logrus.Warn("trying to list without FTP connection. Call Connect first")
		return os.ErrClosed
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entries := exchange.conn.ListObjects(ctx, exchange.config.Bucket, minio.ListObjectsOptions{
		Prefix:    exchange.remote,
		Recursive: true,
	})

	var names []string
	for entry := range entries {
		name := entry.Key
		matched, _ := path.Match(pattern, name)

		if matched && entry.LastModified.Before(before) {
			if err := exchange.conn.RemoveObject(ctx, exchange.config.Bucket, name,
				minio.RemoveObjectOptions{}); err != nil {
				logrus.Errorf("cannot remove %s from %s: %v", name, exchange, err)
			} else {
				logrus.Infof("removed %s from %s", name, exchange)
			}
		}
	}

	logrus.Infof("list from %s: %v", exchange, names)
	return nil
}

func (exchange *S3Exchange) Name() string {
	return exchange.config.Name
}

func (exchange *S3Exchange) String() string {
	return fmt.Sprintf("s3 %s - %s", exchange.config.Endpoint, exchange.config.Bucket)
}
