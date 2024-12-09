package deps

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *client) GetPackage(name string) (*PackageVersions, error) {

	safeName := url.PathEscape(name)

	url := c.baseURL + "/systems/GO/packages/" + safeName
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't make the get request to %q", url)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read all of the body")
	}

	var pkg PackageVersions
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

type Version struct {
	VersionKey VersionKey `json:"versionKey"`
	IsDefault  bool       `json:"isDefault"` // Indicates if it's default version
}

type PackageKey struct {
	System string `json:"system"` // (e.g. GO)
	Name   string `json:"name"`
}

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackageVersions struct {
	PackageKey PackageKey `json:"packageKey"` // Unique identifier for the package
	Versions   []Version  `json:"versions"`   // All available versions
}
