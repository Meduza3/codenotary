package api

import (
	"codenotary/internal"
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

	// Print the project name to the console
	fmt.Printf("Received GET /dependency for project: %s\n", projectName)

	// Retrieve the dependency graph
	fmt.Println("Fetching dependency graph...")
	dependencyGraph, err := internal.Client.GetDependencies(projectName)
	if err != nil {
		http.Error(w, "Failed to fetch dependencies", http.StatusInternalServerError)
		fmt.Printf("Error fetching dependency graph: %v\n", err)
		return
	}
	fmt.Printf("Dependency Graph: %v\n", dependencyGraph)

	// Retrieve all related projects from the dependency graph
	fmt.Println("Fetching all projects from graph...")
	dependenciesProjects, err := internal.Client.GetAllProjectsFromGraph(dependencyGraph)
	if err != nil {
		fmt.Printf("Error fetching related projects: %v\n", err)
	}
	fmt.Printf("Fetched Projects: %v\n", dependenciesProjects)

	// Prepare the structured JSON response
	type Dependency struct {
		ID    string  `json:"id"`
		Score float64 `json:"score"`
	}

	response := struct {
		Message      string       `json:"message"`
		ProjectName  string       `json:"project_name"`
		Dependencies []Dependency `json:"dependencies"`
	}{
		Message:      "Project received successfully",
		ProjectName:  projectName,
		Dependencies: []Dependency{},
	}

	// Populate dependencies
	fmt.Println("Populating dependencies...")
	for _, project := range dependenciesProjects {
		if project == nil {
			fmt.Println("Encountered nil project, skipping...")
			continue
		}
		fmt.Printf("Project: ID=%s, Score=%f\n", project.ProjectKey.ID, project.Scorecard.OverallScore)
		response.Dependencies = append(response.Dependencies, Dependency{
			ID:    project.ProjectKey.ID,
			Score: project.Scorecard.OverallScore,
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
