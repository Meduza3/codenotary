package main

import (
	"codenotary/internal"
	"codenotary/internal/deps"
	"codenotary/internal/sqlite"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var err error
	internal.Db, err = sql.Open("sqlite3", internal.Database)
	if err != nil {
		log.Fatalf("Failed it all! %v", err)
	}
	sqlite.Create(internal.Db)
	internal.Client = deps.NewClient(internal.Db)
	mux := http.NewServeMux()
	mux.HandleFunc("/dependency/", HandleGetDependencies)

	// Start the HTTP server
	port := "8080"
	log.Printf("Server is running on port %s", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
