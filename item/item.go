package item

import (
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

	if priorityOrder[a[i].Status] != priorityOrder[a[j].Status] {
		return priorityOrder[a[i].Status] < priorityOrder[a[j].Status]
	}

	if !a[i].Due.IsZero() && a[j].Due.IsZero() {
		return true
	}
	if a[i].Due.IsZero() && !a[j].Due.IsZero() {
		return false
	}
	if !a[i].Due.IsZero() && !a[j].Due.IsZero() {
		return a[i].Due.Before(a[j].Due)
	}

	return false
}
