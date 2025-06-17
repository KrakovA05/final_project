package api

import (
	"database/sql"
	"net/http"
)

var dbConn *sql.DB

func Init(db *sql.DB) {
	dbConn = db
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r, dbConn)
	default:
		http.NotFound(w, r)
	}
}

//
