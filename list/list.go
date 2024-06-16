package list

import (
	"one-list/item"
	"one-list/storage"
	"sort"
	"time"
)

type List struct {
	storage storage.Storage
	items   []item.Item
	freeID  int
}

func (list *List) Init(storage storage.Storage) {
	list.storage = storage

	list.items = list.storage.FetchAllItems()

	list.freeID = 1
	for _, item := range list.items {
		if item.ID >= list.freeID {
			list.freeID = item.ID + 1
		}
	}
}

func (list *List) GetItems() []item.Item {
	list.updateItems()
	sort.Sort(item.ByProgress(list.items))

	return list.items
}

func (list *List) GetItem(id int) *item.Item {
	_, item := list.findItemByID(id)
	if item != nil {
		return item
	}

	return nil
}

func (list *List) SetItem(item item.Item) {
	index, _ := list.findItemByID(item.ID)
	if index == 0 {
		return
	}

	list.items[index] = item
	list.storage.UpdateItem(item)
}

func (list *List) DeleteItem(id int) {
	index, item := list.findItemByID(id)
	if item == nil {
		return
	}

	list.items[index] = list.items[len(list.items)-1]
	list.items = list.items[:len(list.items)-1]

	list.storage.DeleteItem(id)
}

func (list *List) NewItem(name string) {
	newItem := item.Item{
		ID:                list.freeID,
		Status:            item.NEW,
		Display_status:    item.NOT_STARTED,
		Name:              name,
		Description:       "",
		Due:               time.Time{},
		Reminder_interval: 0,
		Last_update:       time.Time{},
	}
	list.freeID++

	list.items = append(list.items, newItem)

	list.storage.AddItem(newItem)
}

func (list *List) updateItems() {
	for i, _ := range list.items {
		listItem := &((list.items)[i])
		if listItem.Reminder_interval != 0 && listItem.Last_update.AddDate(0, 0, listItem.Reminder_interval).Before(time.Now()) {
			listItem.Status = item.NOTIFICATION
		} else if !listItem.Due.IsZero() && listItem.Due.Before(time.Now()) {
			listItem.Status = item.OVERDUE
		} else if listItem.Last_update.IsZero() {
			listItem.Status = item.NEW
		} else {
			listItem.Status = listItem.Display_status
		}
	}
}

func (list *List) findItemByID(id int) (int, *item.Item) {
	for index, item := range list.items {
		if item.ID == id {
			return index, &item
		}
	}
	return 0, nil
}
