package main

import (
	"final/pkg/api"
	"final/pkg/db"
	"log"
	"net/http"
	"os"
)

const defaultPort = "7540"
const webDir = "./web"

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	api.Init(db.DB)

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Server running at http://localhost:%s/\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed: ", err)
	}
}

//
