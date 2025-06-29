package api

import (
	"database/sql"
	"final/pkg/db"
	"net/http"
	"strconv"
)

type TaskResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request, dbConn *sql.DB) {
	tasks, err := db.Tasks(dbConn, 50)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "database error: " + err.Error()})
		return
	}

	resp := map[string][]map[string]string{
		"tasks": {},
	}

	for _, t := range tasks {
		taskMap := map[string]string{
			"id":      strconv.FormatInt(t.ID, 10),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		}
		resp["tasks"] = append(resp["tasks"], taskMap)
	}

	writeJson(w, resp)
}
