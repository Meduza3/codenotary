package main

import (
	"codenotary/internal/deps"
	"codenotary/internal/sqlite"
	"database/sql"
	"fmt"
	"log"
	"sync"

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

	sqlite.InsertProject(db, *project)

	dependencies, err := client.GetDependencies(name)
	if err != nil {
		log.Fatalf("Failed it all at dependencies! %v", err)
	}

	sqlite.InsertDependencyGraph(db, name, *dependencies)

	var wg sync.WaitGroup

	for _, node := range dependencies.Nodes {
		wg.Add(1)
		go func(node deps.Node) {
			defer wg.Done()
			if node.Relation == "SELF" {
				return
			}
			project, err := client.GetProject(node.VersionKey.Name)
			if err != nil {
				log.Fatalf("Failed it all! %v", err)
			}
			sqlite.InsertProject(db, *project)
			deps, err := client.GetDependencies(node.VersionKey.Name)
			sqlite.InsertDependencyGraph(db, node.VersionKey.Name, *deps)
		}(node)
	}

	wg.Wait()
}
