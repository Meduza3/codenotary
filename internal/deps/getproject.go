package deps

import (
	"codenotary/internal/models"
	"codenotary/internal/sqlite"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

func (c *Client) GetProject(projectKey string) (*models.Project, error) {
	if project, err := sqlite.GetProject(c.db, projectKey); project != nil && err == nil {
		return project, nil
	}

	safeName := url.PathEscape(projectKey)

	url := c.baseURL + "/projects/" + safeName

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

	var project models.Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	//fmt.Printf("Raw response for project %s: %s\n", projectKey, string(body))

	return &project, nil
}

func (c *Client) GetAllProjectsFromGraph(graph *models.DependencyGraph) (succesfulProjects []*models.Project, skipped []string, erro error) {
	var wg sync.WaitGroup

	// Structure to hold results from goroutines
	type ProjectResult struct {
		ProjectName string
		Project     *models.Project
		Err         error
	}

	results := make(chan ProjectResult, len(graph.Nodes)) // Buffer size equals number of nodes

	// Loop through each node in the graph
	for _, node := range graph.Nodes {
		wg.Add(1)
		go func(node models.Node) {
			defer wg.Done()

			// Skip the "SELF" relation
			if node.Relation == "SELF" {
				return
			}

			// Fetch the project
			project, err := c.GetProject(node.VersionKey.Name)
			if err != nil {
				// Log the error for debugging
				// fmt.Printf("GetAllProjectsFromGraph: Error fetching project %s: %v\n", node.VersionKey.Name, err)
				err = fmt.Errorf("failure getting project %q: %w", node.VersionKey.Name, err)
			}

			// Send result back to the channel
			results <- ProjectResult{ProjectName: node.VersionKey.Name, Project: project, Err: err}
		}(node)
	}

	// Close results channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var projects []*models.Project
	var skippedProjects []string
	var errs []error

	// Process results from the channel
	for res := range results {
		if res.Err != nil {
			// Accumulate errors for a combined error report
			skippedProjects = append(skippedProjects, res.ProjectName)
			errs = append(errs, res.Err)
		} else if res.Project != nil {
			// Append valid projects
			projects = append(projects, res.Project)
		} else if res.Project == nil {
			skippedProjects = append(skippedProjects, res.ProjectName)
		}
	}

	// Combine errors if any occurred
	var err error
	if len(errs) > 0 {
		// Format all errors into a single error
		err = fmt.Errorf("multiple errors: %v", errs)
	}

	return projects, skippedProjects, err
}
