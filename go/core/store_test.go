package core

import (
	"log"
	"testing"
)

func TestStoreList(t *testing.T) {

	store := Store{".."}
	list := List(store)
	for _, item := range list {
		log.Printf("Item %s, folder %t", item.path, item.dir)
	}

}
