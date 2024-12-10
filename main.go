package main

import (
	"codenotary/internal/deps"
	"codenotary/internal/sqlite"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const database string = "codenotary.db"

func main() {
	db, err := sql.Open("sqlite3", database)
	sqlite.Create(db)
	client := deps.NewClient(db)

	name := "github.com/cli/cli/v2"

	packageversions, err := client.GetPackage(name)
	if err != nil {
		log.Fatalf("Failed it all! %v", err)
	}

	err = sqlite.StorePackageVersions(db, packageversions)
	if err != nil {
		log.Fatalf("Failed it all! %v", err)
	}

	project, err := client.GetProject(name)
	if err != nil {
		log.Fatalf("Failed it all! %v", err)
	}
	fmt.Println(project.Description)

	sqlite.InsertProject(db, project)

	dependencies, err := client.GetDependencies(name)
	if err != nil {
		log.Fatalf("Failed it all at dependencies! %v", err)
	}

	// sqlite.InsertDependencyGraph(db, name, dependencies)

	projects, err := client.GetAllProjectsFromGraph(dependencies)
	if err != nil {
		fmt.Println(err)
	}

	err = sqlite.InsertProjects(db, projects)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Finished inserting all projects at once")
	}
}
