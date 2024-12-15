package deps

import (
	"codenotary/internal/models"
	"codenotary/internal/sqlite"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) GetPackage(name string) (*models.PackageVersions, error) {
	if pkg, err := sqlite.GetPackageVersions(c.db, name); pkg != nil && err == nil {
		return pkg, nil
	}

	safeName := url.PathEscape(name)

	url := c.baseURL + "/systems/GO/packages/" + safeName
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't make the get request to %q: %w", url, err)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read all of the body")
	}

	var pkg models.PackageVersions
	if err := json.Unmarshal(body, &pkg); err != nil {
		return nil, fmt.Errorf("couldn't parse JSON: %v", err)
	}

	return &pkg, nil
}

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
