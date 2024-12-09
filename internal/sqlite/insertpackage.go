package sqlite

import (
	"codenotary/internal/deps"
	"database/sql"
	"fmt"
)

// Insert or update package and its versions
func StorePackageVersions(db *sql.DB, pv *deps.PackageVersions) error {
	// Insert or ignore into packages table
	_, err := db.Exec(`
			INSERT OR IGNORE INTO packages (system, name) VALUES (?, ?)`,
		pv.PackageKey.System, pv.PackageKey.Name)
	if err != nil {
		return fmt.Errorf("failed to insert package: %v", err)
	}

	// For each version, insert or replace
	for _, v := range pv.Versions {
		isDefault := 0
		if v.IsDefault {
			isDefault = 1
		}
		_, err = db.Exec(`
					INSERT OR REPLACE INTO package_versions (system, name, version, is_default)
					VALUES (?, ?, ?, ?)`,
			v.VersionKey.System, v.VersionKey.Name, v.VersionKey.Version, isDefault)
		if err != nil {
			return fmt.Errorf("failed to insert version: %v", err)
		}
	}
	return nil
}

// Example function to get the latest (default) version of a project
func GetLatestVersion(db *sql.DB, system, name string) (string, error) {
	var version string
	err := db.QueryRow(`
		SELECT version FROM package_versions
		WHERE system = ? AND name = ? AND is_default = 1
		LIMIT 1`, system, name).Scan(&version)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("no default version found for %s/%s", system, name)
	} else if err != nil {
		return "", fmt.Errorf("failed to get latest version: %v", err)
	}
	return version, nil
}
