package frontend

import "one-list/list"

type Frontend interface {
	Init(st list.List, username string, password string)
	Close()
}
