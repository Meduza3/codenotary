package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

func InsertDependencyGraph(db *sql.DB, projectID string, graph *models.DependencyGraph) error {
	if graph == nil {
		return fmt.Errorf("nil graph")
	}

	
	for idx, node := range graph.Nodes {
		errors := strings.Join(node.Errors, ";") 

		query := `
			INSERT INTO dependency_nodes (project_id, graph_id, node_index, system, name, version, bundled, relation, errors)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
		args := []interface{}{
			projectID, 
			projectID, idx, node.VersionKey.System, node.VersionKey.Name, node.VersionKey.Version,
			node.Bundled, node.Relation, errors,
		}

		

		_, err := db.Exec(query, args...)
		if err != nil {
			log.Printf("Error inserting node at index %d: %v", idx, err)
			return fmt.Errorf("failed to insert node at index %d: %v", idx, err)
		}
	}

	
	for _, edge := range graph.Edges {
		query := `
			INSERT INTO dependency_edges (project_id, graph_id, from_node_index, to_node_index, requirement)
			VALUES (?, ?, ?, ?, ?)`
		args := []interface{}{
			projectID, 
			projectID, edge.FromNode, edge.ToNode, edge.Requirement,
		}

		

		_, err := db.Exec(query, args...)
		if err != nil {
			log.Printf("Error inserting edge from %d to %d: %v", edge.FromNode, edge.ToNode, err)
			return fmt.Errorf("failed to insert edge from %d to %d: %v", edge.FromNode, edge.ToNode, err)
		}
	}

	log.Println("Dependency graph inserted successfully.")
	return nil
}
