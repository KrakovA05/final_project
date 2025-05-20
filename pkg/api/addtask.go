package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final/pkg/db"
	"io"
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
	today := now.Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("неверный формат даты")
	}

	if task.Repeat != "" {
		if !t.After(today) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("неверный формат правила повторения")
			}
			task.Date = next
		}
	} else {
		// ВАЖНО: возможно убрать проверку на дату в прошлом и не возвращать ошибку
		// или заменить её на предупреждение, если в тестах так нужно.
		if !t.After(today) {
			// вместо ошибки просто выставляем дату на сегодня или на завтра?
			// или убираем это условие, если так тесты хотят
			return errors.New("дата задачи в прошлом")
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": "ошибка чтения тела запроса"})
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &task); err != nil {
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
		writeJson(w, map[string]string{"error": "ошибка при добавлении задачи в базу"})
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
