package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

func AddOrUpdateDependency(db *sql.DB, projectID string, dep models.Node) error {
	query := `
		INSERT INTO dependency_nodes (
			project_id, graph_id, node_index, system, name, version, bundled, relation, errors, ossf_score
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(project_id, graph_id, node_index) DO UPDATE SET
			system=excluded.system,
			name=excluded.name,
			version=excluded.version,
			bundled=excluded.bundled,
			relation=excluded.relation,
			errors=excluded.errors,
			ossf_score=excluded.ossf_score;
	`
	ossfScore, err := CalculateOpenSSF(db, dep.VersionKey.Name) 
	if err != nil {
		ossfScore = -1
	}
	_, err = db.Exec(query,
		projectID,
		projectID,           
		dep.VersionKey.Name, 
		dep.VersionKey.System,
		dep.VersionKey.Name,
		dep.VersionKey.Version,
		dep.Bundled,
		dep.Relation,
		strings.Join(dep.Errors, ";"),
		ossfScore,
	)
	if err != nil {
		return fmt.Errorf("failed to add or update dependency: %v", err)
	}
	return nil
}

func CalculateOpenSSF(db *sql.DB, depName string) (float64, error) {
	query := `SELECT scorecard_overall_score FROM project WHERE id = ? LIMIT 1`
	var score float64
	err := db.QueryRow(query, depName).Scan(&score)
	if err != nil {
		if err == sql.ErrNoRows {
			
			
			return 0.0, nil
		}
		return 0.0, fmt.Errorf("failed to retrieve OpenSSF score: %v", err)
	}
	return score, nil
}


func GetDependency(db *sql.DB, projectID, depName string) (*models.Node, error) {
	query := `
		SELECT system, name, version, bundled, relation, errors, ossf_score
		FROM dependency_nodes
		WHERE project_id = ? AND name = ?
		LIMIT 1;
	`
	row := db.QueryRow(query, projectID, depName)

	var dep models.Node
	var errorsStr string
	var ossfScore float64
	if err := row.Scan(&dep.VersionKey.System, &dep.VersionKey.Name, &dep.VersionKey.Version, &dep.Bundled, &dep.Relation, &errorsStr, &ossfScore); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, fmt.Errorf("failed to get dependency: %v", err)
	}

	dep.Errors = strings.Split(errorsStr, ";")
	
	

	return &dep, nil
}


func DeleteDependency(db *sql.DB, projectID, depName string) error {
	query := `
		DELETE FROM dependency_nodes
		WHERE project_id = ? AND name = ?;
	`
	result, err := db.Exec(query, projectID, depName)
	if err != nil {
		return fmt.Errorf("failed to delete dependency: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no dependency found to delete")
	}
	return nil
}


func ListDependencies(db *sql.DB, name string, minScore float64) ([]models.Node, error) {
	query := `
		SELECT system, name, version, bundled, relation, errors, ossf_score
		FROM dependency_nodes
		WHERE 1=1
	`
	args := []interface{}{}

	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}

	if minScore > 0 {
		query += " AND ossf_score >= ?"
		args = append(args, minScore)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list dependencies: %v", err)
	}
	defer rows.Close()

	var deps []models.Node
	for rows.Next() {
		var dep models.Node
		var errorsStr string
		var ossfScore float64
		if err := rows.Scan(&dep.VersionKey.System, &dep.VersionKey.Name, &dep.VersionKey.Version, &dep.Bundled, &dep.Relation, &errorsStr, &ossfScore); err != nil {
			return nil, fmt.Errorf("failed to scan dependency: %v", err)
		}
		dep.Errors = strings.Split(errorsStr, ";")
		
		deps = append(deps, dep)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating dependencies: %v", err)
	}

	return deps, nil
}
