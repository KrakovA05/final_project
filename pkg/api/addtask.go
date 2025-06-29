package api

import (
	"database/sql"
	"encoding/json"
	"final/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = today.Format("20060102")
		return nil
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	next, err := NextDate(today, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	if afterNow(today, t) {
		if len(task.Repeat) == 0 {
			task.Date = today.Format("20060102")
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(dbConn, &task)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка базы данных: " + err.Error()})
		return
	}

	writeJson(w, map[string]int64{"id": id})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(dbConn, id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJson(w, map[string]string{
		"id":      id,
		"date":    task.Date,
		"title":   task.Title,
		"comment": task.Comment,
		"repeat":  task.Repeat,
	})
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	type inputTask struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}

	var input inputTask
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJson(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if input.ID == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"})
		return
	}
	idInt, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{"error": "неверный формат id"})
		return
	}

	if input.Title == "" {
		writeJson(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	task := db.Task{
		ID:      idInt,
		Date:    input.Date,
		Title:   input.Title,
		Comment: input.Comment,
		Repeat:  input.Repeat,
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(dbConn, &task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	err := db.DeleteTask(dbConn, id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{}) // успешный ответ — пустой JSON
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "не указан идентификатор"})
		return
	}

	task, err := db.GetTask(dbConn, id)
	if err != nil {
		writeJson(w, map[string]string{"error": "задача не найдена"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(dbConn, id); err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, map[string]string{})
		return
	}

	// Парсим текущую дату задачи
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "неверный формат даты"})
		return
	}

	next, err := NextDate(date, task.Date, task.Repeat)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(dbConn, next, id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}
