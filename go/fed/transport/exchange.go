package transport

import "time"

type UpdatesCh chan []string
type Exchange interface {
	ID() string

	/** Connect to the transport using local as folder where download/upload files are */
	Connect(remoteRoot, localRoot string) (UpdatesCh, error)
	Disconnect()

	List(since time.Time) ([]string, error)

	/** Pull the requested file (full local) from the transport to the storage folder */
	Push(file string) (int64, error)

	/** Pull the requested location from the transport to the storage folder*/
	Pull(loc string) (int64, error)

	/** Delete all files matching pattern and older than specified time */
	Delete(pattern string, before time.Time) error

	Name() string
	String() string
}
