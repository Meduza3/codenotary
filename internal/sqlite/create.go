package sqlite

import (
	"database/sql"
	"fmt"
)

func Create(db *sql.DB) error {
	// Create the project table
	projectTable := `
			CREATE TABLE IF NOT EXISTS project (
					id TEXT PRIMARY KEY,
					open_issues_count INTEGER,
					stars_count INTEGER,
					forks_count INTEGER,
					license TEXT,
					description TEXT,
					homepage TEXT,
					scorecard_date TEXT,
					scorecard_repo_name TEXT,
					scorecard_repo_commit TEXT,
					scorecard_version TEXT,
					scorecard_commit TEXT,
					scorecard_overall_score REAL
			);
			`
	_, err := db.Exec(projectTable)
	if err != nil {
		return fmt.Errorf("failed to create project table: %v", err)
	}

	// Create a table for scorecard checks
	checksTable := `
			CREATE TABLE IF NOT EXISTS scorecard_checks (
					project_id TEXT,
					name TEXT,
					short_description TEXT,
					url TEXT,
					score REAL,
					reason TEXT,
					details TEXT,
					FOREIGN KEY (project_id) REFERENCES project(id)
			);
			`
	_, err = db.Exec(checksTable)
	if err != nil {
		return fmt.Errorf("failed to create scorecard_checks table: %v", err)
	}

	packagesTable := `
	CREATE TABLE IF NOT EXISTS packages (
		system TEXT,
		name TEXT,
		PRIMARY KEY (system, name)
	);
	`
	if _, err := db.Exec(packagesTable); err != nil {
		return fmt.Errorf("failed to create packages table: %v", err)
	}

	versionsTable := `
	CREATE TABLE IF NOT EXISTS package_versions (
		system TEXT,
		name TEXT,
		version TEXT,
		is_default INTEGER,
		PRIMARY KEY (system, name, version),
		FOREIGN KEY (system, name) REFERENCES packages(system, name)
	);
	`
	if _, err := db.Exec(versionsTable); err != nil {
		return fmt.Errorf("failed to create package_versions table: %v", err)
	}

	dependency_nodes := `
	CREATE TABLE IF NOT EXISTS dependency_nodes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id TEXT,				
  graph_id TEXT,           -- to reference which graph this node belongs to, if needed
  node_index INTEGER,         -- index of this node in the original response
  system TEXT,
  name TEXT,
  version TEXT,
  bundled BOOLEAN,
  relation TEXT,              -- SELF, DIRECT, INDIRECT
  errors TEXT                 -- could store as JSON or newline-delimited
	);
	`

	if _, err := db.Exec(dependency_nodes); err != nil {
		return fmt.Errorf("failed to create package_versions table: %v", err)
	}

	dependency_edges := `
	CREATE TABLE IF NOT EXISTS dependency_edges (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id TEXT,
  graph_id TEXT,
  from_node_index INTEGER,
  to_node_index INTEGER,
  requirement TEXT
	);
	`
	if _, err := db.Exec(dependency_edges); err != nil {
		return fmt.Errorf("failed to create package_versions table: %v", err)
	}

	return nil
}
