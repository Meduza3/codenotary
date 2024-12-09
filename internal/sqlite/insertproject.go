package sqlite

import (
	"codenotary/internal/deps"
	"database/sql"
	"fmt"
	"strings"
)

// insertProject inserts the project and related checks into the database.
func InsertProject(db *sql.DB, p deps.Project) error {
	// Insert into project table
	insertProjectSQL := `
	INSERT INTO project (
			id, open_issues_count, stars_count, forks_count, license, description, homepage,
			scorecard_date, scorecard_repo_name, scorecard_repo_commit, scorecard_version,
			scorecard_commit, scorecard_overall_score
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);
	`
	_, err := db.Exec(insertProjectSQL,
		p.ProjectKey.ID,
		p.OpenIssuesCount,
		p.StarsCount,
		p.ForksCount,
		p.License,
		p.Description,
		p.Homepage,
		p.Scorecard.Date,
		p.Scorecard.Repository.Name,
		p.Scorecard.Repository.Commit,
		p.Scorecard.Scorecard.Version,
		p.Scorecard.Scorecard.Commit,
		p.Scorecard.OverallScore,
	)
	if err != nil {
		return fmt.Errorf("failed to insert project: %v", err)
	}

	// Insert checks
	insertCheckSQL := `
	INSERT INTO scorecard_checks (
			project_id, name, short_description, url, score, reason, details
	) VALUES (?,?,?,?,?,?,?);
	`
	for _, check := range p.Scorecard.Checks {
		detailsStr := strings.Join(check.Details, "\n")
		_, err := db.Exec(insertCheckSQL,
			p.ProjectKey.ID,
			check.Name,
			check.Documentation.ShortDescription,
			check.Documentation.URL,
			check.Score,
			check.Reason,
			detailsStr,
		)
		if err != nil {
			return fmt.Errorf("failed to insert scorecard check: %v", err)
		}
	}

	return nil
}
