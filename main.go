package main

import (
	"one-list/frontend"
	"one-list/frontend/web"
	"one-list/list"
	"one-list/storage"
	"one-list/storage/sqlite"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const USERNAME = "a"
const PASSWORD = "b"

func main() {
	var init sync.WaitGroup
	init.Add(1)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	var storage storage.Storage = &sqlite.Sqlite{}
	go storage.Init(&init)
	defer storage.Close()
	init.Wait()

	var list list.List = list.List{}
	list.Init(storage)

	var frontend frontend.Frontend = &web.Web{}
	go frontend.Init(list, USERNAME, PASSWORD)
	defer frontend.Close()

	<-c
}
