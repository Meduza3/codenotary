package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
)



func GetPackageVersions(db *sql.DB, name string) (*models.PackageVersions, error) {
	
	system := "GO"
	var pkgName string
	err := db.QueryRow(`
		SELECT name FROM packages
		WHERE system = ? AND name = ?`,
		system, name).Scan(&pkgName)
	if err == sql.ErrNoRows {
		return nil, nil 
	} else if err != nil {
		return nil, fmt.Errorf("error querying packages: %v", err)
	}

	
	rows, err := db.Query(`
		SELECT version, is_default FROM package_versions
		WHERE system = ? AND name = ?`,
		system, name)
	if err != nil {
		return nil, fmt.Errorf("error querying package_versions: %v", err)
	}
	defer rows.Close()

	var versions []models.Version
	for rows.Next() {
		var version string
		var isDefault int
		if err := rows.Scan(&version, &isDefault); err != nil {
			return nil, fmt.Errorf("error scanning package_versions row: %v", err)
		}
		v := models.Version{
			VersionKey: models.VersionKey{
				System:  system,
				Name:    name,
				Version: version,
			},
			IsDefault: isDefault == 1,
		}
		versions = append(versions, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating package_versions rows: %v", err)
	}

	
	if len(versions) == 0 {
		return nil, nil
	}

	pv := &models.PackageVersions{
		PackageKey: models.PackageKey{
			System: system,
			Name:   name,
		},
		Versions: versions,
	}

	return pv, nil
}
