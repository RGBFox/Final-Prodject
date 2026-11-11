package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
		return arg, fmt.Errorf("такой задачи не существует")
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
		return http.StatusBadRequest, err
	}
	if n == 0 {
		return http.StatusBadRequest, sql.ErrNoRows
	}
	return http.StatusOK, nil
}
