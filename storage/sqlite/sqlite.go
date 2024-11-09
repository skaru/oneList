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

type Sqlite struct {
	db *sql.DB
}

func (sql Sqlite) DeleteItem(ID int) {
	stmt, err := sql.db.Prepare("DELETE FROM " + DB_TABLE + " WHERE id = ?")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		log.Println(err)
	}
}

func (sql Sqlite) AddItem(item item.Item) {
	stmt, err := sql.db.Prepare("INSERT INTO " + DB_TABLE + " ('name', display_status, reminder_interval) VALUES (?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(item.Name, item.Display_status, item.Reminder_interval)
	if err != nil {
		log.Println(err)
	}
}

func (sql Sqlite) UpdateItem(item item.Item) {
	stmt, err := sql.db.Prepare("UPDATE " + DB_TABLE + " SET Display_status = ?, name = ?, description = ?, due = ?, reminder_interval = ?, last_update = ? WHERE id = ?")
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

	_, err = stmt.Exec(
		item.Display_status,
		item.Name,
		item.Description,
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

	var item item.Item
	var due int64
	var lastUpdate int64
	err = stmt.QueryRow(ID).Scan(
		&item.ID,
		&item.Display_status,
		&item.Name,
		&item.Description,
		&due,
		&item.Reminder_interval,
		&lastUpdate)
	if err != nil {
		log.Println(err)
	}

	item = sql.parseDates(item, due, lastUpdate)
	return item
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
		var item item.Item
		var due int64
		var lastUpdate int64

		err := rows.Scan(
			&item.ID,
			&item.Display_status,
			&item.Name,
			&item.Description,
			&due,
			&item.Reminder_interval,
			&lastUpdate)
		if err != nil {
			log.Println(err)
		}

		item = sql.parseDates(item, due, lastUpdate)
		items = append(items, item)
	}

	return items
}

func (sql Sqlite) parseDates(item item.Item, due int64, lastUpdate int64) item.Item {
	if due != 0 {
		item.Due = time.Unix(due, 0)
	} else {
		item.Due = time.Time{}
	}

	if lastUpdate != 0 {
		item.Last_update = time.Unix(lastUpdate, 0)
	} else {
		item.Last_update = time.Time{}
	}

	return item
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

	_, err = sq.db.Query("CREATE TABLE `" + DB_TABLE + "` (`id` INTEGER PRIMARY KEY NOT NULL, `display_status` INTEGER NOT NULL DEFAULT '1', `name` VARCHAR(255) NOT NULL, `description` TEXT(65535) NOT NULL DEFAULT '', `due` INTEGER NOT NULL DEFAULT '0', `reminder_interval` INTEGER NOT NULL DEFAULT '0', `last_update` INTEGER NOT NULL DEFAULT '0')")
	if err != nil {
		log.Println(err)
	}

	init.Done()
}
