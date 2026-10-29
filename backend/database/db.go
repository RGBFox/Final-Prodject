package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func Init() error {
	dbFile := "scheduler.db"
	path := os.Getenv("TODO_DBFILE")
	if path != "" {
		dbFile = path
	}

	_, err := os.Stat(dbFile)

	var install bool
	if os.IsNotExist(err) {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()
	if install {
		res, err := db.Exec(`
	CREATE TABLE scheduler (
   		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date CHAR(8) NOT NULL DEFAULT '',
		title TEXT NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat TEXT NOT NULL DEFAULT ''
    );
	CREATE INDEX date_index 
	ON scheduler (date);
	`)
		if err != nil {
			return err
		}
		log.Print("база данных создана", res)

	}
	return nil
}
