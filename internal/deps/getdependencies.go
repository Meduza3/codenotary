package deps

import (
	"codenotary/internal/models"
	"codenotary/internal/sqlite"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// First checks the database for a dependency graph, then fetches from api
func (c *Client) GetDependencies(name string) (*models.DependencyGraph, error) {
	graph, err := sqlite.GetDependencyGraph(c.db, name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving dependency graph from DB: %v", err)
	}

	if graph != nil {
		log.Printf("Dependency graph for project %q found in the database.", name)
		return graph, nil
	}

	safeName := url.PathEscape(name)

	latestVersion, err := c.GetLatestVersionByProjectId(name)
	if err != nil {
		return nil, fmt.Errorf("Couldn't get latest version of project ID %q: %v", name, err)
	}
	url := fmt.Sprintf("%s/systems/GO/packages/%s/versions/%s:dependencies", c.baseURL, safeName, latestVersion)
	resp, err := http.Get(url)
	//fmt.Println(url)
	if err != nil {
		return nil, fmt.Errorf("Couldn't make the get request to %q", url)
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read all of the body")
	}
	var dependencyGraph models.DependencyGraph
	if err := json.Unmarshal(body, &dependencyGraph); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	err = sqlite.InsertDependencyGraph(c.db, name, &dependencyGraph)
	if err != nil {
		//handle error
	}
	return &dependencyGraph, nil
}
