package main

import (
	"codenotary/internal"
	"codenotary/internal/models"
	"codenotary/internal/sqlite"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func HandleGetDependencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the project name by trimming the prefix
	encodedProjectName := strings.TrimPrefix(r.URL.Path, "/dependency/")

	// Decode the URL-encoded project name
	projectName, err := url.QueryUnescape(encodedProjectName)
	if err != nil || projectName == "" {
		http.Error(w, "Invalid or missing project name", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received GET /dependency for project: %s\n", projectName)

	//store it in DB
	project, err := internal.Client.GetProject(projectName)
	if err != nil {
		//handle error!
	}
	err = sqlite.InsertProject(internal.Db, project)
	if err != nil {
		//handle error!
	}

	// Retrieve the dependency graph
	fmt.Print("Fetching dependency graph...")
	dependencyGraph, err := internal.Client.GetDependencies(projectName)
	if err != nil {
		http.Error(w, "Failed to fetch dependencies", http.StatusInternalServerError)
		fmt.Printf("Error fetching dependency graph: %v\n", err)
		return
	}
	// Retrieve all related projects from the dependency graph
	dependenciesProjects, skipped, err := internal.Client.GetAllProjectsFromGraph(dependencyGraph)
	if err != nil {
		fmt.Printf("Error fetching related projects: %v\n", err)
	}

	for _, skippedProject := range skipped {
		proj := emptyProjectFromName(skippedProject)
		dependenciesProjects = append(dependenciesProjects, &proj)
	}

	// Prepare the structured JSON response
	type Dependency struct {
		ID          string         `json:"id"`
		Score       float64        `json:"score"`
		CheckScores map[string]int `json:"check_scores,omitempty"`
	}

	mainScores, err := sqlite.GetScoresByProjectID(internal.Db, projectName)
	if err != nil {
		//handle error!
	}
	response := struct {
		MainScores   map[string]int `json:"main_scores"`
		Message      string         `json:"message"`
		ProjectName  string         `json:"project_name"`
		Dependencies []Dependency   `json:"dependencies"`
	}{
		Message:      "No dependencies =) Hiring Marcin is a great idea",
		MainScores:   mainScores,
		ProjectName:  projectName,
		Dependencies: []Dependency{},
	}

	// Populate dependencies
	fmt.Println("Populating dependencies...")
	sqlite.InsertProjects(internal.Db, dependenciesProjects)
	for _, project := range dependenciesProjects {
		if project == nil {
			fmt.Println("Encountered nil project, skipping...")
			continue
		}
		checkScores, err := sqlite.GetScoresByProjectID(internal.Db, project.ProjectKey.ID)
		if err != nil {
			fmt.Printf("Error fetching scores for project %s: %v\n", project.ProjectKey.ID, err)
			checkScores = nil
		}
		//fmt.Printf("Project: ID=%s, Score=%f\n", project.ProjectKey.ID, project.Scorecard.OverallScore)
		response.Dependencies = append(response.Dependencies, Dependency{
			ID:          project.ProjectKey.ID,
			Score:       project.Scorecard.OverallScore,
			CheckScores: checkScores,
		})
	}

	// Print the JSON response to the console
	fmt.Printf("Generated JSON Response: %s\n", toJSONString(response))

	// Respond with the JSON structure
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Printf("Error encoding response: %v\n", err)
	}
}

// Helper function to convert any struct to a JSON string
func toJSONString(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func emptyProjectFromName(name string) models.Project {
	// Create the project key, ID, pathname
	// and init an empty project
	project := models.Project{
		ProjectKey: models.ProjectKey{
			ID: name,
		},
		Scorecard: models.Scorecard{
			OverallScore: -1,
		},
	}

	// Return empty project
	return project
}
