package main

import (
	"codenotary/internal"
	"codenotary/internal/models"
	"codenotary/internal/sqlite"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)


type AddOrUpdateDependencyRequest struct {
	ProjectName string      `json:"project_name"`
	Dependency  models.Node `json:"dependency"`
}


func HandleAddOrUpdateDependency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddOrUpdateDependencyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.ProjectName == "" || req.Dependency.VersionKey.Name == "" {
		http.Error(w, "Project name and dependency name are required", http.StatusBadRequest)
		return
	}

	
	err := sqlite.AddOrUpdateDependency(internal.Db, req.ProjectName, req.Dependency)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add/update dependency: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dependency added/updated successfully"))
}


func HandleGetDependency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	
	trimmedPath := strings.TrimPrefix(r.URL.Path, "/dependency/get/")
	if trimmedPath == r.URL.Path {
		
		http.Error(w, "Invalid URL format. Use /dependency/get/{projectName}/{dependencyName}", http.StatusBadRequest)
		return
	}

	parts := strings.SplitN(trimmedPath, "/", 2)
	if len(parts) != 2 {
		http.Error(w, "Invalid URL format. Use /dependency/get/{projectName}/{dependencyName}", http.StatusBadRequest)
		return
	}
	projectName, depName := parts[0], parts[1]

	dep, err := sqlite.GetDependency(internal.Db, projectName, depName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving dependency: %v", err), http.StatusInternalServerError)
		return
	}

	if dep == nil {
		http.Error(w, "Dependency not found", http.StatusNotFound)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dep)
}


func HandleDeleteDependency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/dependency/"), "/")
	if len(parts) != 2 {
		http.Error(w, "Invalid URL format. Use /dependency/{projectName}/{dependencyName}", http.StatusBadRequest)
		return
	}
	projectName, depName := parts[0], parts[1]

	err := sqlite.DeleteDependency(internal.Db, projectName, depName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting dependency: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Dependency deleted successfully"))
}


func HandleListDependencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	
	name := r.URL.Query().Get("name")
	minScoreStr := r.URL.Query().Get("min_score")
	minScore := 0.0
	var err error
	if minScoreStr != "" {
		minScore, err = strconv.ParseFloat(minScoreStr, 64)
		if err != nil {
			http.Error(w, "Invalid min_score value", http.StatusBadRequest)
			return
		}
	}

	deps, err := sqlite.ListDependencies(internal.Db, name, minScore)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing dependencies: %v", err), http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deps)
}
