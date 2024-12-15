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
	
	if err != nil {
		return nil, fmt.Errorf("Couldn't make the get request to %q: %w", url, err)
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

	

	return &project, nil
}

func (c *Client) GetAllProjectsFromGraph(graph *models.DependencyGraph) (succesfulProjects []*models.Project, skipped []string, erro error) {
	var wg sync.WaitGroup

	
	type ProjectResult struct {
		ProjectName string
		Project     *models.Project
		Err         error
	}

	results := make(chan ProjectResult, len(graph.Nodes)) 

	
	for _, node := range graph.Nodes {
		wg.Add(1)
		go func(node models.Node) {
			defer wg.Done()

			
			if node.Relation == "SELF" {
				return
			}

			
			project, err := c.GetProject(node.VersionKey.Name)
			if err != nil {
				
				
				err = fmt.Errorf("failure getting project %q: %w", node.VersionKey.Name, err)
			}

			
			results <- ProjectResult{ProjectName: node.VersionKey.Name, Project: project, Err: err}
		}(node)
	}

	
	go func() {
		wg.Wait()
		close(results)
	}()

	var projects []*models.Project
	var skippedProjects []string
	var errs []error

	
	for res := range results {
		if res.Err != nil {
			
			skippedProjects = append(skippedProjects, res.ProjectName)
			errs = append(errs, res.Err)
		} else if res.Project != nil {
			
			projects = append(projects, res.Project)
		} else if res.Project == nil {
			skippedProjects = append(skippedProjects, res.ProjectName)
		}
	}

	
	var err error
	if len(errs) > 0 {
		
		err = fmt.Errorf("multiple errors: %v", errs)
	}

	return projects, skippedProjects, err
}
