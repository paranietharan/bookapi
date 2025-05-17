package main

import (
	"bookapi/pkg/router"
	"bookapi/pkg/store"
	"time"
)

func main() {
	store.StartBookCleanup(40 * time.Second)

	router.StartServer()
}
