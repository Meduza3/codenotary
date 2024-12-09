package sqlite

import (
	"codenotary/internal/deps"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func InsertDependencyGraph(db *sql.DB, projectID string, graph deps.DependencyGraph) error {
	// Insert nodes into dependency_nodes
	for idx, node := range graph.Nodes {
		errors := strings.Join(node.Errors, ";") // Combine errors into a single string

		query := `
			INSERT INTO dependency_nodes (project_id, graph_id, node_index, system, name, version, bundled, relation, errors)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		args := []interface{}{
			projectID, // Use project ID as graph_id (or modify this if you implement a separate graph table)
			projectID, idx, node.VersionKey.System, node.VersionKey.Name, node.VersionKey.Version,
			node.Bundled, node.Relation, errors,
		}

		log.Printf("Executing query: %s with args: %v", query, args)

		_, err := db.Exec(query, args...)
		if err != nil {
			log.Printf("Error inserting node at index %d: %v", idx, err)
			return fmt.Errorf("failed to insert node at index %d: %v", idx, err)
		}
	}

	// Insert edges into dependency_edges
	for _, edge := range graph.Edges {
		query := `
			INSERT INTO dependency_edges (project_id, graph_id, from_node_index, to_node_index, requirement)
			VALUES (?, ?, ?, ?, ?)`
		args := []interface{}{
			projectID, // Use project ID as graph_id (or modify this if you implement a separate graph table)
			projectID, edge.FromNode, edge.ToNode, edge.Requirement,
		}

		log.Printf("Executing query: %s with args: %v", query, args)

		_, err := db.Exec(query, args...)
		if err != nil {
			log.Printf("Error inserting edge from %d to %d: %v", edge.FromNode, edge.ToNode, err)
			return fmt.Errorf("failed to insert edge from %d to %d: %v", edge.FromNode, edge.ToNode, err)
		}
	}

	log.Println("Dependency graph inserted successfully.")
	return nil
}
