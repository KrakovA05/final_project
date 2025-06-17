package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"final/pkg/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	date, err := time.Parse("20060102", task.Date)
	if err != nil {
		return errors.New("дата представлена в формате, отличном от 20060102")
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if date.Before(today) && task.Repeat != "" {
		newDate, err := shiftDate(date, task.Repeat)
		if err != nil {
			return errors.New("правило повторения указано в неправильном формате")
		}
		task.Date = newDate.Format("20060102")

		date = newDate
	}

	if date.Before(today) {
		return errors.New("дата не может быть в прошлом")
	}

	return nil
}

func shiftDate(from time.Time, repeat string) (time.Time, error) {
	parts := strings.Fields(repeat)
	if len(parts) != 2 {
		return time.Time{}, errors.New("неправильный формат правила повторения")
	}

	unit := parts[0]
	n, err := strconv.Atoi(parts[1])
	if err != nil || n <= 0 {
		return time.Time{}, errors.New("неправильное число в правиле повторения")
	}

	next := from
	today := time.Now()
	for !next.After(today) {
		switch unit {
		case "d":
			next = next.AddDate(0, 0, n)
		case "w":
			next = next.AddDate(0, 0, 7*n)
		case "m":
			next = next.AddDate(0, n, 0)
		case "y":
			next = next.AddDate(n, 0, 0)
		default:
			return time.Time{}, errors.New("неподдерживаемая единица повтора")
		}
	}

	return next, nil
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

//
