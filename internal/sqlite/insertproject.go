package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
	"strings"
)


func InsertProject(db *sql.DB, p *models.Project) error {
	if p == nil {
		return fmt.Errorf("Nil project")
	}
	
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

func InsertProjects(db *sql.DB, projects []*models.Project) error {
	
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	
	insertProjectStmt, err := tx.Prepare(`
		INSERT INTO project (
			id, open_issues_count, stars_count, forks_count, license, description, homepage,
			scorecard_date, scorecard_repo_name, scorecard_repo_commit, scorecard_version,
			scorecard_commit, scorecard_overall_score
		) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare project insert statement: %w", err)
	}
	defer insertProjectStmt.Close()

	
	insertCheckStmt, err := tx.Prepare(`
		INSERT INTO scorecard_checks (
			project_id, name, short_description, url, score, reason, details
		) VALUES (?,?,?,?,?,?,?)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare scorecard checks insert statement: %w", err)
	}
	defer insertCheckStmt.Close()

	var errs []error

	for _, p := range projects {
		if p == nil {
			continue
		}

		
		_, perr := insertProjectStmt.Exec(
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
		if perr != nil {
			
			errs = append(errs, fmt.Errorf("failed to insert project (ID: %s): %w", p.ProjectKey.ID, perr))
			
			continue
		}

		
		for _, check := range p.Scorecard.Checks {
			detailsStr := strings.Join(check.Details, "\n")
			_, cerr := insertCheckStmt.Exec(
				p.ProjectKey.ID,
				check.Name,
				check.Documentation.ShortDescription,
				check.Documentation.URL,
				check.Score,
				check.Reason,
				detailsStr,
			)
			if cerr != nil {
				
				errs = append(errs, fmt.Errorf("failed to insert scorecard check for project (ID: %s, check: %s): %w",
					p.ProjectKey.ID, check.Name, cerr))
				
			}
		}
	}

	
	commitErr := tx.Commit()
	if commitErr != nil {
		return fmt.Errorf("failed to commit transaction: %w", commitErr)
	}

	
	if len(errs) > 0 {
		
		var sb strings.Builder
		sb.WriteString("Some inserts failed:\n")
		for _, e := range errs {
			sb.WriteString(e.Error())
			sb.WriteString("\n")
		}
		return fmt.Errorf(sb.String())
	}

	return nil
}
