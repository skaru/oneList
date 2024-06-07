package item

import (
	"sort"
	"time"
)

type Progress int

const (
	ON_HOLD Progress = iota
	NOT_STARTED
	IN_PROGRESS
	DONE
	NEW
	OVERDUE
	NOTIFICATION
)

type Item struct {
	ID                int
	Status            Progress
	Display_status    Progress
	Name              string
	Description       string
	Due               time.Time
	Reminder_interval int
	Last_update       time.Time
}

func NewItem(id int, name string) Item {
	return Item{
		ID:                id,
		Status:            NEW,
		Display_status:    NOT_STARTED,
		Name:              name,
		Description:       "",
		Due:               time.Time{},
		Reminder_interval: 0,
		Last_update:       time.Time{},
	}
}

func UpdateAndSortItems(items []Item) {
	updateItems(&items)
	sort.Sort(ByProgress(items))
}

func updateItems(items *[]Item) {
	for i, _ := range *items {
		item := &((*items)[i])
		if item.Reminder_interval != 0 && item.Last_update.AddDate(0, 0, item.Reminder_interval).Before(time.Now()) {
			item.Status = NOTIFICATION
		} else if !item.Due.IsZero() && item.Due.Before(time.Now()) {
			item.Status = OVERDUE
		} else if item.Last_update.IsZero() {
			item.Status = NEW
		} else {
			item.Status = item.Display_status
		}
	}
}

type ByProgress []Item

func (a ByProgress) Len() int      { return len(a) }
func (a ByProgress) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByProgress) Less(i, j int) bool {
	priorityOrder := map[Progress]int{
		NOTIFICATION: 0,
		OVERDUE:      1,
		NEW:          2,
		IN_PROGRESS:  3,
		ON_HOLD:      4,
		NOT_STARTED:  5,
		DONE:         6,
	}
	return priorityOrder[a[i].Status] < priorityOrder[a[j].Status]
}
