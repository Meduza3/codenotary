package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
	"strings"
)

// GetProject retrieves a Project from the database based on the projectKey.
// It returns nil if the project is not found.
func GetProject(db *sql.DB, projectKey string) (*models.Project, error) {
	// Retrieve the project details
	var p models.Project
	var scorecardDate, repoName, repoCommit, scorecardVersion, scorecardCommit string
	var overallScore float64

	projectQuery := `
		SELECT 
			id, open_issues_count, stars_count, forks_count, license, description, homepage,
			scorecard_date, scorecard_repo_name, scorecard_repo_commit, scorecard_version,
			scorecard_commit, scorecard_overall_score
		FROM project
		WHERE id = ?
	`

	err := db.QueryRow(projectQuery, projectKey).Scan(
		&p.ProjectKey.ID,
		&p.OpenIssuesCount,
		&p.StarsCount,
		&p.ForksCount,
		&p.License,
		&p.Description,
		&p.Homepage,
		&scorecardDate,
		&repoName,
		&repoCommit,
		&scorecardVersion,
		&scorecardCommit,
		&overallScore,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Project not found
	} else if err != nil {
		return nil, fmt.Errorf("error querying project: %v", err)
	}

	p.Scorecard.Date = scorecardDate
	p.Scorecard.Repository.Name = repoName
	p.Scorecard.Repository.Commit = repoCommit
	p.Scorecard.Scorecard.Version = scorecardVersion
	p.Scorecard.Scorecard.Commit = scorecardCommit
	p.Scorecard.OverallScore = overallScore

	// Retrieve associated scorecard checks
	checksQuery := `
		SELECT 
			name, short_description, url, score, reason, details
		FROM scorecard_checks
		WHERE project_id = ?
	`

	rows, err := db.Query(checksQuery, projectKey)
	if err != nil {
		return nil, fmt.Errorf("error querying scorecard_checks: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var check models.ScorecardCheck
		var detailsStr string

		if err := rows.Scan(
			&check.Name,
			&check.Documentation.ShortDescription,
			&check.Documentation.URL,
			&check.Score,
			&check.Reason,
			&detailsStr,
		); err != nil {
			return nil, fmt.Errorf("error scanning scorecard_checks row: %v", err)
		}

		check.Details = strings.Split(detailsStr, "\n")
		p.Scorecard.Checks = append(p.Scorecard.Checks, check)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scorecard_checks rows: %v", err)
	}

	return &p, nil
}
