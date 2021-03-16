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
	Push(file string) error

	/** Pull the requested location from the transport to the storage folder*/
	Pull(loc string) error

	String() string
}
