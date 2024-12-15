package main

import (
	"codenotary/internal"
	"codenotary/internal/deps"
	"codenotary/internal/sqlite"
	"database/sql"
	"fmt"
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
	err = sqlite.Create(internal.Db)
	if err != nil {
		fmt.Println(err)
	}
	internal.Client = deps.NewClient(internal.Db)
	mux := http.NewServeMux()

	internal.Client.GetProject("github.com/cli/cli")
	mux.HandleFunc("/dependency/", HandleGetDependencies)

	mux.HandleFunc("/dependency/add", HandleAddOrUpdateDependency) // POST
	mux.HandleFunc("/dependency/get/", HandleGetDependency)        // GET /dependency/get/{projectName}/{depName}
	mux.HandleFunc("/dependency/delete/", HandleDeleteDependency)  // DELETE /dependency/delete/{projectName}/{depName}
	mux.HandleFunc("/dependencies", HandleListDependencies)        // GET /dependencies?name=xyz&min_score=50

	port := "8080"
	log.Printf("Server is running on port %s", port)
	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
