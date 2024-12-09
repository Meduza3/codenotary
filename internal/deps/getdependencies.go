package deps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *client) GetDependencies(name string) (*DependencyGraph, error) {
	safeName := url.PathEscape(name)

	latestVersion, err := c.GetLatestVersionByProjectId(name)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get latest version of project ID %q: %v", name, err)
	}
	url := fmt.Sprintf("%s/systems/GO/packages/%s/versions/%s:dependencies", c.baseURL, safeName, latestVersion)
	resp, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't make the get request to %q", url)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read all of the body")
	}
	var project DependencyGraph
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	return &project, nil
}

type Node struct {
	VersionKey VersionKey `json:"versionKey"`
	Bundled    bool       `json:"bundled"`
	Relation   string     `json:"relation"`
	Errors     []string   `json:"errors"`
}

type Edge struct {
	FromNode    int    `json:"fromNode"`
	ToNode      int    `json:"toNode"`
	Requirement string `json:"requirement"`
}

type DependencyGraph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
	Error string `json:"error"`
}
