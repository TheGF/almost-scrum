package library

import (
	"almost-scrum/core"
	"crypto/sha256"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type ExtendedAttr struct {
	Hash  []byte `json:"hash"`
	Owner string `json:"owner"`
}

type ExtendedAttrMap struct {
	Version    string                   `json:"version"`
	Entries    map[string]*ExtendedAttr `json:"entries"`
	lastAccess time.Time
	dirty      int32
}

const attrsFileName = ".ash-xAttrs.json"

var cache = make(map[string]*ExtendedAttrMap)
var cacheC = make(chan string)
var cacheWg sync.WaitGroup

func cacheSync() {
	for path := range cacheC {
		extendedAttrMap, found := cache[path]
		if found {
			if atomic.SwapInt32(&extendedAttrMap.dirty, 0) == 1 {
				err := core.WriteJSON(filepath.Join(path, attrsFileName), extendedAttrMap)
				if err != nil {
					logrus.Errorf("Cannot save extended attrs to %s: %v", path, err)
				} else {
//					logrus.Debugf("xAttr index %s updated successfully", path)
				}
			}
		} else {
//			logrus.Debugf("cannot find path %s in extended attributes cache", path)
		}
		cacheWg.Done()
	}
}

func init() {
	go cacheSync()
}

func getSHA256(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h256 := sha256.New()
	if _, err := io.Copy(h256, file); err != nil {
		return nil, err
	}
	return h256.Sum(nil), nil
}

func getExtendedAttrMap(path string) (*ExtendedAttrMap, error) {
	extendedAttrMap, found := cache[path]
	if found {
		extendedAttrMap.lastAccess = time.Now()
		return extendedAttrMap, nil
	}
	p := filepath.Join(path, attrsFileName)
	if _, err := os.Stat(p); os.IsNotExist(err) {
//		logrus.Debugf("xAttrs index %s does not exist. Return empty instance", p)
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
	err := core.ReadJSON(p, extendedAttrMap)
	if err != nil {
		return nil, err
	}
//	logrus.Debugf("successfully read xAttrs index %s, len(entries)=%d", p, len(extendedAttrMap.Entries))
	extendedAttrMap.lastAccess = time.Now()
	cache[path] = extendedAttrMap
	return extendedAttrMap, nil
}

func saveExtendedAttrMapLater(path string, extendedAttrMap *ExtendedAttrMap) {
	atomic.StoreInt32(&extendedAttrMap.dirty, 1)
	cacheWg.Add(1)
	cacheC <- path
}

func getExtendedAttr(dir string, name string) (*ExtendedAttr, error) {
	extendedAttrMap, err := getExtendedAttrMap(dir)
	if err != nil {
		return nil, err
	}
	extendedAttr, found := extendedAttrMap.Entries[name]
	if found {
		return extendedAttr, nil
	}

	path := filepath.Join(dir, name)
	hash, err := getSHA256(path)
	if err != nil {
		return nil, err
	}

	owner, err := getFileOwner(path)
	if err != nil {
		return nil, err
	}

//	logrus.Debugf("no xAttrs for file %s. Setting new xAttrs with hash=%s and owner=%s", path, hash, owner)
	extendedAttr = &ExtendedAttr{
		Hash:  hash,
		Owner: owner,
	}
	extendedAttrMap.Entries[name] = extendedAttr
	saveExtendedAttrMapLater(dir, extendedAttrMap)
	return extendedAttr, nil
}

func setExtendedAttr(dir string, name string, extendedAttr *ExtendedAttr) error {
	extendedAttrMap, err := getExtendedAttrMap(dir)
	if err != nil {
		return err
	}
	if extendedAttr == nil {
		delete(extendedAttrMap.Entries, name)
		logrus.Debugf("Deleted xAttrs for %s/%s", dir, name)
	} else {
		extendedAttrMap.Entries[name] = extendedAttr
		logrus.Debugf("Updated xAttrs for %s/%s to hash=%v, owner=%s", dir, name,
			extendedAttr.Hash, extendedAttr.Owner)
	}
	saveExtendedAttrMapLater(dir, extendedAttrMap)
	return nil
}

func setOwner(dir string, name string, owner string) error {
	extendedAttrMap, err := getExtendedAttrMap(dir)
	if err != nil {
		return err
	}

	extendedAttr, found :=  extendedAttrMap.Entries[name]
	if found {
		extendedAttr.Owner = owner
	} else {
		path := filepath.Join(dir, name)
		hash, err := getSHA256(path)
		if err != nil {
			return err
		}
		extendedAttrMap.Entries[name] = &ExtendedAttr{
			Hash:  hash,
			Owner: owner,
		}
	}
	saveExtendedAttrMapLater(dir, extendedAttrMap)
	return nil
}
