package main

import (
	"bookapi/pkg/router"
	"bookapi/pkg/store"
	"time"
)

func main() {
	store.StartBookCleanup(1 * time.Minute)

	router.StartServer()
}
