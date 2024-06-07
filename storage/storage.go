package storage

import (
	"one-list/item"
	"sync"
)

type Storage interface {
	AddItem(item item.Item)
	DeleteItem(ID int)
	UpdateItem(item item.Item)
	FetchAllItems() []item.Item
	FetchItem(ID int) item.Item
	Init(init *sync.WaitGroup)
	Close()
}
