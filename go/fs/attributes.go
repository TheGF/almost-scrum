package fs

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"
	"time"
)

type ExtendedAttr struct {
	Owner      string    `json:"owner"`
	ImportHash []byte    `json:"importHash"`
	ExportHash []byte    `json:"exportHash"`
	Public     bool      `json:"public_"`
	Modified   time.Time `json:"modified"`
}

type ExtendedAttrMap struct {
	Version    string                   `json:"version"`
	Entries    map[string]*ExtendedAttr `json:"entries"`
	lastAccess time.Time
	dirty      int32
	lock       sync.RWMutex
}

const AttrsFileName = ".ash-xAttrs.json"

var cache = make(map[string]*ExtendedAttrMap)
var cacheC = make(chan string)
var cacheWg sync.WaitGroup

func cacheSync() {
	for path := range cacheC {
		extendedAttrMap, found := cache[path]
		if found {
			extendedAttrMap.lock.RLock()
			if atomic.SwapInt32(&extendedAttrMap.dirty, 0) == 1 {
				err := WriteJSON(filepath.Join(path, AttrsFileName), extendedAttrMap)
				if err != nil {
					logrus.Errorf("Cannot save extended attrs to %s: %v", path, err)
				}
			}
			extendedAttrMap.lock.RUnlock()
		}
		cacheWg.Done()
	}
}

func init() {
	go cacheSync()
}

func getExtendedAttrMap(path string) (*ExtendedAttrMap, error) {
	extendedAttrMap, found := cache[path]
	if found {
		extendedAttrMap.lastAccess = time.Now()
		return extendedAttrMap, nil
	}
	p := filepath.Join(path, AttrsFileName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		extendedAttrMap = &ExtendedAttrMap{
			Version:    "1.0",
			Entries:    make(map[string]*ExtendedAttr),
			lastAccess: time.Time{},
		}
		cache[path] = extendedAttrMap
		return extendedAttrMap, nil
	} else if err != nil {
		logrus.Errorf("cannot get xAttrs %s: %v", path, err)
		return nil, err
	}
	extendedAttrMap = &ExtendedAttrMap{}
	err := ReadJSON(p, extendedAttrMap)
	if err != nil {
		return nil, err
	}
	extendedAttrMap.lastAccess = time.Now()
	cache[path] = extendedAttrMap
	return extendedAttrMap, nil
}

func saveExtendedAttrMapLater(path string, extendedAttrMap *ExtendedAttrMap) {
	atomic.StoreInt32(&extendedAttrMap.dirty, 1)
	cacheWg.Add(1)
	cacheC <- path
}

func GetExtendedAttr(path string) (*ExtendedAttr, error) {
	dir, name := filepath.Split(path)
	extendedAttrMap, err := getExtendedAttrMap(dir)
	if err != nil {
		return nil, err
	}
	extendedAttr, found := extendedAttrMap.Entries[name]
	if found {
		return extendedAttr, nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	owner, err := GetFileOwner(path)
	if err != nil {
		return nil, err
	}

	extendedAttr = &ExtendedAttr{
		Owner: owner,
	}
	extendedAttrMap.lock.Lock()
	extendedAttrMap.Entries[name] = extendedAttr
	extendedAttrMap.lock.Unlock()

	saveExtendedAttrMapLater(dir, extendedAttrMap)
	return extendedAttr, nil
}

func SetExtendedAttr(path string, extendedAttr *ExtendedAttr) error {
	dir, name := filepath.Split(path)
	extendedAttrMap, err := getExtendedAttrMap(dir)
	if err != nil {
		return err
	}
	extendedAttrMap.lock.Lock()
	if extendedAttr == nil {
		delete(extendedAttrMap.Entries, name)
		logrus.Infof("Modified xAttrs for %s", path)
	} else {
		extendedAttrMap.Entries[name] = extendedAttr
		logrus.Infof("updated attr for %s: %#v", path, extendedAttr)
	}
	extendedAttrMap.lock.Unlock()
	saveExtendedAttrMapLater(dir, extendedAttrMap)
	return nil
}

var versionsRegex = regexp.MustCompile(`(.*?)(~(\d+\.)+\d+)?(\.\w*)?$`)

func ParsePath(path string) (dir string, prefix string, version string, ext string, err error) {
	dir = filepath.Dir(path)
	name := filepath.Base(path)
	match := versionsRegex.FindStringSubmatch(name)
	if len(match) != 5 {
		err = os.ErrInvalid
		logrus.Errorf("Cannot parse %s: %v", path, err)
		return
	}

	prefix = match[1]
	version = match[2]
	ext = match[4]
	err = nil

	return
}
