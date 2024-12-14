package sqlite

import (
	"codenotary/internal/models"
	"database/sql"
	"fmt"
)

// GetPackageVersions retrieves the PackageVersions for a given system and name from the database.
// It returns nil if the package is not found.
func GetPackageVersions(db *sql.DB, name string) (*models.PackageVersions, error) {
	// Retrieve package
	system := "GO"
	var pkgName string
	err := db.QueryRow(`
		SELECT name FROM packages
		WHERE system = ? AND name = ?`,
		system, name).Scan(&pkgName)
	if err == sql.ErrNoRows {
		return nil, nil // Package not found
	} else if err != nil {
		return nil, fmt.Errorf("error querying packages: %v", err)
	}

	// Retrieve versions
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

	// If no versions found, treat as package not found
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
