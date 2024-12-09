package deps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
