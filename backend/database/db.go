package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// Init выбирает файл SQLite и создаёт таблицу и индекс, если файла ещё нет.
// Возвращает ошибку и путь к файлу БД.
func Init() (error, string) {
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

	// Временный пул нужен для подготовки БД; рабочий пул создаётся в main.
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err, ""
	}
	defer db.Close()
	if install {
		// Индекс по date предназначен для поиска и сортировки задач по дате.
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
			return err, ""
		}
		log.Print("база данных создана", res)

	}
	return nil, dbFile
}
