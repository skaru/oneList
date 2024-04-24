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
)

type Item struct {
	ID                int
	Status            Progress
	Name              string
	Description       string
	Priority          int
	Due               time.Time
	Reminder_interval int
	Last_update       time.Time
	Creation_date     time.Time
}

func Sort(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		// Place items with Status 0 at the bottom
		if items[i].Status == 0 && items[j].Status != 0 {
			return false
		} else if items[i].Status != 0 && items[j].Status == 0 {
			return true
		}

		// Check if reminder interval is due
		iReminderDue := items[i].Reminder_interval > 0 && items[i].Last_update.AddDate(0, 0, items[i].Reminder_interval).Before(time.Now())
		jReminderDue := items[j].Reminder_interval > 0 && items[j].Last_update.AddDate(0, 0, items[j].Reminder_interval).Before(time.Now())

		if iReminderDue != jReminderDue {
			return iReminderDue
		}

		// Check for Last_update
		iUpdateIsZero := items[i].Last_update.IsZero()
		jUpdateIsZero := items[j].Last_update.IsZero()

		if iUpdateIsZero != jUpdateIsZero {
			return iUpdateIsZero
		} else if iUpdateIsZero {
			return items[i].Creation_date.Before(items[j].Creation_date)
		}

		// Sort by Status in descending order
		if items[i].Status != items[j].Status {
			return items[i].Status > items[j].Status
		}

		// For items with the same Status, check for due dates
		iDueIsZero := items[i].Due.IsZero()
		jDueIsZero := items[j].Due.IsZero()

		if iDueIsZero != jDueIsZero {
			return jDueIsZero
		}

		return items[i].Due.Before(items[j].Due)
	})
}
