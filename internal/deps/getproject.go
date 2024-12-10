package deps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

func (c *client) GetProject(projectKey string) (*Project, error) {
	safeName := url.PathEscape(projectKey)

	url := c.baseURL + "/projects/" + safeName

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

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	return &project, nil
}

func (c *client) GetProjectAndAllDependencies(projectKey string) (*result, error) {
	safeName := url.PathEscape(projectKey)

	url := c.baseURL + "/projects/" + safeName

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

	var project Project
	if err := json.Unmarshal(body, &project); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}

	return nil, nil
}

func (c *client) GetAllProjectsFromGraph(graph *DependencyGraph) ([]*Project, error) {
	var wg sync.WaitGroup
	type ProjectResult struct {
		Project *Project
		Err     error
	}
	results := make(chan ProjectResult, len(graph.Nodes))
	for _, node := range graph.Nodes {
		wg.Add(1)
		go func(node Node) {
			defer wg.Done()
			if node.Relation == "SELF" {
				results <- ProjectResult{Project: nil, Err: nil}
				return
			}
			project, err := c.GetProject(node.VersionKey.Name)
			if err != nil {
				err = fmt.Errorf("failure getting project %q: %w", node.VersionKey.Name, err)
			}
			results <- ProjectResult{Project: project, Err: err}
		}(node)
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	var projects []*Project
	var errs []error
	for res := range results {
		if res.Err != nil {
			errs = append(errs, res.Err)
		} else if res.Project != nil {
			projects = append(projects, res.Project)
		}
	}
	var err error
	if len(errs) > 0 {
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
