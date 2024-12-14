package main

import (
	"codenotary/api"
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
	mux.HandleFunc("/dependency/", api.HandleGetDependencies)
	port := "8080"
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		panic(err)
	}
}
