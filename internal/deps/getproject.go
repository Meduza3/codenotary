package deps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

func (c *Client) GetProject(projectKey string) (*Project, error) {
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

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	//fmt.Printf("Raw response for project %s: %s\n", projectKey, string(body))

	return &project, nil
}

func (c *Client) GetProjectAndAllDependencies(projectKey string) (*result, error) {
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

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	return nil, nil
}

func (c *Client) GetAllProjectsFromGraph(graph *DependencyGraph) ([]*Project, error) {
	var wg sync.WaitGroup

	// Structure to hold results from goroutines
	type ProjectResult struct {
		Project *Project
		Err     error
	}

	results := make(chan ProjectResult, len(graph.Nodes)) // Buffer size equals number of nodes

	// Loop through each node in the graph
	for _, node := range graph.Nodes {
		wg.Add(1)
		go func(node Node) {
			defer wg.Done()

			// Skip the "SELF" relation
			if node.Relation == "SELF" {
				results <- ProjectResult{Project: nil, Err: nil}
				return
			}

			// Fetch the project
			project, err := c.GetProject(node.VersionKey.Name)
			if err != nil {
				// Log the error for debugging
				fmt.Printf("GetAllProjectsFromGraph: Error fetching project %s: %v\n", node.VersionKey.Name, err)
				err = fmt.Errorf("failure getting project %q: %w", node.VersionKey.Name, err)
			}

			// Send result back to the channel
			results <- ProjectResult{Project: project, Err: err}
		}(node)
	}

	// Close results channel after all goroutines finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var projects []*Project
	var errs []error

	// Process results from the channel
	for res := range results {
		if res.Err != nil {
			// Accumulate errors for a combined error report
			errs = append(errs, res.Err)
		} else if res.Project != nil {
			// Append valid projects
			projects = append(projects, res.Project)
		}
	}

	// Combine errors if any occurred
	var err error
	if len(errs) > 0 {
		// Format all errors into a single error
		err = fmt.Errorf("multiple errors: %v", errs)
	}

	return projects, err
}

type result struct{}

type Project struct {
	ProjectKey      ProjectKey `json:"projectKey"`
	OpenIssuesCount int        `json:"openIssuesCount"`
	StarsCount      int        `json:"starsCount"`
	ForksCount      int        `json:"forksCount"`
	License         string     `json:"license"`
	Description     string     `json:"description"`
	Homepage        string     `json:"homepage"`
	Scorecard       Scorecard  `json:"scorecard"`
}

type ProjectKey struct {
	ID string `json:"id"`
}

type Scorecard struct {
	Date         string           `json:"date"`
	Repository   Repository       `json:"repository"`
	Scorecard    ScorecardInfo    `json:"scorecard"`
	Checks       []ScorecardCheck `json:"checks"`
	OverallScore float64          `json:"overallScore"`
	Metadata     []string         `json:"metadata"`
}

type Repository struct {
	Name   string `json:"name"`
	Commit string `json:"commit"`
}

type ScorecardInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

type ScorecardCheck struct {
	Name          string             `json:"name"`
	Documentation CheckDocumentation `json:"documentation"`
	Score         float64            `json:"score"`
	Reason        string             `json:"reason"`
	Details       []string           `json:"details"`
}

type CheckDocumentation struct {
	ShortDescription string `json:"shortDescription"`
	URL              string `json:"url"`
}
