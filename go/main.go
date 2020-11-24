package main

import "almost-scrum/core"

func main() {
	store := core.Store{Path: "."}
	core.ListStore(store)
}
