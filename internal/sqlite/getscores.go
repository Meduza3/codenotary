package sqlite

import (
	"database/sql"
	"fmt"
)

func GetScoresByProjectID(db *sql.DB, projectID string) (map[string]int, error) {
	// Query to fetch name and score for the given project ID
	query := `SELECT name, score FROM scorecard_checks WHERE project_id = ?`

	// Prepare the statement to avoid SQL injection
	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Map to hold the results
	result := make(map[string]int)

	// Iterate through the rows
	for rows.Next() {
		var name string
		var score int
		if err := rows.Scan(&name, &score); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result[name] = score
	}

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating through rows: %v", err)
	}

	return result, nil
}
