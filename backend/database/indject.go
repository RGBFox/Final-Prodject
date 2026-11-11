package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

// deleteTas удаляет задачу по ID  и сообщает, если запись не найдена.
func DeleteTask(db *sql.DB, id int) (int, error) {
	res, err := db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		log.Printf("ошибка удаления задач: %v\n", err)
		return http.StatusInternalServerError, fmt.Errorf("Ошибка удаления задачи")
	}
	n, err := res.RowsAffected()
	if err != nil {
		log.Printf("ошибка удаления задач: %v\n", err)
		return http.StatusInternalServerError, fmt.Errorf("Ошибка поиска задачи")
	}
	if n == 0 {
		return http.StatusBadRequest, fmt.Errorf("задача не найдена")
	}
	return http.StatusOK, nil
}

// SearchByIDTask возвращает одну задачу по параметру id.
func SearchByIDTask(db *sql.DB, id int) (Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`
	res := db.QueryRow(query, sql.Named("id", id))
	var arg Task
	err := res.Scan(&arg.ID, &arg.Date, &arg.Title, &arg.Comment, &arg.Repeat)
	if err != nil {
		return arg, err
	}
	return arg, nil
}

// AddTask сохраняет подготовленную задачу и возвращает ID, назначенный SQLite.
func AddTask(task *Task, db *sql.DB) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat)
	          VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query, sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// updateTask записывает поля задачи и возвращает sql.ErrNoRows, если её ID не найден.
func UpdateTask(task *Task, db *sql.DB) (int, error) {

	query := `UPDATE scheduler 
	SET date = :date,
	  title= :title,
	  comment= :comment,
	  repeat=  :repeat 
	WHERE id = :id;
	`
	res, err := db.Exec(query, sql.Named("id", task.ID), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// UPDATE без совпадений не считается ошибкой SQL, поэтому проверяем число строк.
	n, err := res.RowsAffected()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if n == 0 {
		return http.StatusBadRequest, sql.ErrNoRows
	}
	return http.StatusOK, nil
}

// SearchByIDTask возвращает одну задачи по параметрам или общий список задач.
func SearchTask(db *sql.DB, search string) ([]Task, int, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50`
	var args []any
	// Инициализированный пустой срез кодируется в JSON как [], а не null.
	tasks := []Task{}

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			// поиск по дате
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date ORDER BY date LIMIT 50`
			// Дату из пользовательского формата приводим к формату хранения в SQLite.
			args = append(args, sql.Named("date", t.Format("20060102")))
		} else {
			// поиск по тексту
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :pattern OR comment LIKE :pattern ORDER BY date LIMIT 50`
			// Знаки % вокруг параметра позволяют искать совпадение внутри текста.
			args = append(args, sql.Named("pattern", "%"+search+"%"))
		}
	}

	res, err := db.Query(query, args...)
	if err != nil {
		log.Println(err)
		return tasks, http.StatusInternalServerError, fmt.Errorf("Ошибка базы данных")
	}
	defer res.Close()
	for res.Next() {
		var arg Task
		err = res.Scan(&arg.ID, &arg.Date, &arg.Title, &arg.Comment, &arg.Repeat)
		if err != nil {
			log.Println(err)
			return tasks, http.StatusInternalServerError, fmt.Errorf("такой задачи не существует")
		}
		tasks = append(tasks, arg)
	}
	// Next может завершиться из-за ошибки чтения; проверяем её перед отправкой ответа.
	if err := res.Err(); err != nil {
		log.Printf("ошибки чтения: %v\n", err)
		return tasks, http.StatusInternalServerError, fmt.Errorf("данные не полностью считались")
	}
	return tasks, http.StatusOK, nil

}
