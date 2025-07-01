package api

import (
	"database/sql"
	"net/http"
)

var dbConn *sql.DB

func Init(db *sql.DB) {
	dbConn = db

	http.HandleFunc("/api/nextdate", nextDateHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			addTaskHandler(w, r, dbConn)
		case http.MethodGet:
			getTaskHandler(w, r, dbConn)
		case http.MethodPut:
			updateTaskHandler(w, r, dbConn)
		case http.MethodDelete:
			deleteTaskHandler(w, r, dbConn)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}) // строка 27
		}
	})

	http.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			doneTaskHandler(w, r, dbConn)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}) // строка 36
		}
	})

	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tasksHandler(w, r, dbConn)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"}) // строка 46
		}
	})

}
