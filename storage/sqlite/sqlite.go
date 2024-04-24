package sqlite

import (
	"database/sql"
	"log"
	"one-list/item"
	"os"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

const DB_NAME = "main.db"
const DB_DRIVER = "sqlite"
const DB_TABLE = "items"
const DB_TIME_FORMAT = "2006-01-02 15:04:05"

type sqlItem struct {
	ID                int
	Status            item.Progress
	Name              string
	Description       string
	Priority          int
	Due               int64
	Reminder_interval int
	Last_update       int64
	Creation_date     int64
}

type Sqlite struct {
	db *sql.DB
}

func (sql Sqlite) DeleteItem(ID int) {
	stmt, err := sql.db.Prepare("DELETE FROM " + DB_TABLE + " WHERE id = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Query(ID)
	if err != nil {
		log.Println(err)
	}
}

func (sql Sqlite) CreateItem(name string) {
	stmt, err := sql.db.Prepare("INSERT INTO " + DB_TABLE + " ('name', creation_date) VALUES (?, strftime('%s', 'now'))")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Query(name)
	if err != nil {
		log.Println(err)
	}
}

func (sql Sqlite) UpdateItem(item item.Item) {
	stmt, err := sql.db.Prepare("UPDATE " + DB_TABLE + " SET status = ?, name = ?, description = ?, priority = ?, due = ?, reminder_interval = ?, last_update = ? WHERE id = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	var due int64 = 0
	var lastUpdate int64 = 0

	if !item.Due.IsZero() {
		due = item.Due.Unix()
	}

	if !item.Last_update.IsZero() {
		lastUpdate = item.Last_update.Unix()
	}

	_, err = stmt.Query(
		item.Status,
		item.Name,
		item.Description,
		item.Priority,
		due,
		item.Reminder_interval,
		lastUpdate,
		item.ID,
	)
	if err != nil {
		log.Println(err)
	}
}

func (sql Sqlite) FetchItem(ID int) item.Item {
	stmt, err := sql.db.Prepare("SELECT * FROM " + DB_TABLE + " WHERE id = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	var sqlItem sqlItem
	err = stmt.QueryRow(ID).Scan(
		&sqlItem.ID,
		&sqlItem.Status,
		&sqlItem.Name,
		&sqlItem.Description,
		&sqlItem.Priority,
		&sqlItem.Due,
		&sqlItem.Reminder_interval,
		&sqlItem.Last_update,
		&sqlItem.Creation_date)
	if err != nil {
		log.Println(err)
	}

	var dueDate time.Time
	var lastUpdate time.Time

	if sqlItem.Due != 0 {
		dueDate = time.Unix(sqlItem.Due, 0)
	}

	if sqlItem.Last_update != 0 {
		lastUpdate = time.Unix(sqlItem.Last_update, 0)
	}

	outputItem := item.Item{
		ID:                sqlItem.ID,
		Status:            sqlItem.Status,
		Name:              sqlItem.Name,
		Description:       sqlItem.Description,
		Priority:          sqlItem.Priority,
		Due:               dueDate,
		Reminder_interval: sqlItem.Reminder_interval,
		Last_update:       lastUpdate,
		Creation_date:     time.Unix(sqlItem.Creation_date, 0),
	}

	return outputItem
}

func (sql Sqlite) FetchAllItems() []item.Item {
	stmt, err := sql.db.Prepare("SELECT * FROM " + DB_TABLE)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var items []item.Item
	for rows.Next() {
		var sqlItem sqlItem
		err := rows.Scan(
			&sqlItem.ID,
			&sqlItem.Status,
			&sqlItem.Name,
			&sqlItem.Description,
			&sqlItem.Priority,
			&sqlItem.Due,
			&sqlItem.Reminder_interval,
			&sqlItem.Last_update,
			&sqlItem.Creation_date)
		if err != nil {
			log.Println(err)
		}

		var dueDate time.Time
		var lastUpdate time.Time

		if sqlItem.Due != 0 {
			dueDate = time.Unix(sqlItem.Due, 0)
		}

		if sqlItem.Last_update != 0 {
			lastUpdate = time.Unix(sqlItem.Last_update, 0)
		}

		outputItem := item.Item{
			ID:                sqlItem.ID,
			Status:            sqlItem.Status,
			Name:              sqlItem.Name,
			Description:       sqlItem.Description,
			Priority:          sqlItem.Priority,
			Due:               dueDate,
			Reminder_interval: sqlItem.Reminder_interval,
			Last_update:       lastUpdate,
			Creation_date:     time.Unix(sqlItem.Creation_date, 0),
		}
		items = append(items, outputItem)
	}

	return items
}

func (sql Sqlite) Close() {
	sql.db.Close()
}

func (sq *Sqlite) Init(init *sync.WaitGroup) {
	var err error
	if _, err := os.Stat(DB_NAME); os.IsNotExist(err) {
		file, err := os.Create(DB_NAME)
		if err != nil {
			log.Fatal("Error creating Database:", err)
		}
		defer file.Close()

		log.Println("Database created:", DB_NAME)
	} else {
		log.Println("Database already exists:", DB_NAME)
	}

	sq.db, err = sql.Open(DB_DRIVER, DB_NAME+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	_, err = sq.db.Query("CREATE TABLE `" + DB_TABLE + "` (`id` INTEGER PRIMARY KEY NOT NULL, `status` INTEGER NOT NULL DEFAULT '1', `name` VARCHAR(255) NOT NULL, `description` TEXT(65535) NOT NULL DEFAULT '', `priority` INTEGER NOT NULL DEFAULT '0', `due` INTEGER NOT NULL DEFAULT '0', `reminder_interval` INTEGER NOT NULL DEFAULT '0', `last_update` INTEGER NOT NULL DEFAULT '0', `creation_date` INTEGER NOT NULL)")
	if err != nil {
		log.Println(err)
	}

	init.Done()
}
