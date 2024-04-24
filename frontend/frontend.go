package frontend

import "one-list/storage"

type Frontend interface {
	Init(st storage.Storage, username string, password string)
	Close()
}
