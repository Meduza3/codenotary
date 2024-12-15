package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
	"strings"
)



func GetDependencyGraph(db *sql.DB, projectID string) (*models.DependencyGraph, error) {
	
	nodeRows, err := db.Query(`
			SELECT node_index, system, name, version, bundled, relation, errors
			FROM dependency_nodes
			WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("error querying dependency_nodes: %v", err)
	}
	defer nodeRows.Close()

	var nodes []models.Node
	for nodeRows.Next() {
		var nodeIndex int
		var system, name, version string
		var bundled bool
		var relation string
		var errorsStr sql.NullString

		if err := nodeRows.Scan(&nodeIndex, &system, &name, &version, &bundled, &relation, &errorsStr); err != nil {
			return nil, fmt.Errorf("error scanning dependency_nodes row: %v", err)
		}

		node := models.Node{
			VersionKey: models.VersionKey{
				System:  system,
				Name:    name,
				Version: version,
			},
			Bundled:  bundled,
			Relation: relation,
		}

		if errorsStr.Valid && errorsStr.String != "" {
			node.Errors = strings.Split(errorsStr.String, ";")
		} else {
			node.Errors = []string{}
		}

		nodes = append(nodes, node)
	}

	if err := nodeRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating dependency_nodes rows: %v", err)
	}

	
	if len(nodes) == 0 {
		return nil, nil
	}

	
	edgeRows, err := db.Query(`
			SELECT from_node_index, to_node_index, requirement
			FROM dependency_edges
			WHERE project_id = ?`, projectID)
	if err != nil {
		return nil, fmt.Errorf("error querying dependency_edges: %v", err)
	}
	defer edgeRows.Close()

	var edges []models.Edge
	for edgeRows.Next() {
		var fromNode, toNode int
		var requirement string

		if err := edgeRows.Scan(&fromNode, &toNode, &requirement); err != nil {
			return nil, fmt.Errorf("error scanning dependency_edges row: %v", err)
		}

		edge := models.Edge{
			FromNode:    fromNode,
			ToNode:      toNode,
			Requirement: requirement,
		}

		edges = append(edges, edge)
	}

	if err := edgeRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating dependency_edges rows: %v", err)
	}

	graph := &models.DependencyGraph{
		Nodes: nodes,
		Edges: edges,
	}

	return graph, nil
}
