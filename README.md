# CodeNotary Technical challenge

# Setup
first run ```docker-compose up --build``` to build

then ```docker-compose up``` to run it another time

Visit the frontend at ```localhost:3000/index```
You can query the api at ```localhost:8080```

Supported features:
  - First query local SQLite, then if not found ask the deps API
  - Bar graph of the count of scores
  - Click on any dependency in the list to search among its dependencies
  - A back button to go to the previous project
  - Fast loading time due to asynchronous goroutines

# SQLite schema

### project
id TEXT : Unique project identifier (Primary Key)  
open_issues_count INTEGER : Number of open issues  
stars_count INTEGER : Star count  
forks_count INTEGER : Fork count  
license TEXT : Project license type  
description TEXT : Brief project description  
homepage TEXT : Project homepage URL  
scorecard_date TEXT : Date of the scorecard  
scorecard_repo_name TEXT : Repository name on the scorecard  
scorecard_repo_commit TEXT : Repository commit hash  
scorecard_version TEXT : Scorecard tool version  
scorecard_commit TEXT : Commit hash of the scorecard tool  
scorecard_overall_score REAL : Overall security score  

### scorecard_checks
project_id TEXT : Related project ID (Foreign Key to `project.id`)  
name TEXT : Name of the scorecard check  
short_description TEXT : Summary of the check  
url TEXT : URL for additional details  
score REAL : Check score  
reason TEXT : Reasoning for the score  
details TEXT : Additional details  

### packages
system TEXT : Package ecosystem (e.g., npm, pip)  
name TEXT : Package name  

### package_versions
system TEXT : Related ecosystem (Foreign Key to `packages.system`)  
name TEXT : Related package name (Foreign Key to `packages.name`)  
version TEXT : Package version  
is_default INTEGER : Whether this version is the default (1 for true, 0 for false)  

### dependency_nodes
id INTEGER : Unique node identifier (Primary Key)  
project_id TEXT : Associated project (Foreign Key to `project.id`)  
graph_id TEXT : Identifier for the dependency graph  
node_index INTEGER : Node's index within the graph  
system TEXT : Package ecosystem  
name TEXT : Package name  
version TEXT : Package version  
bundled BOOLEAN : Whether the dependency is bundled  
relation TEXT : Dependency type (`SELF`, `DIRECT`, `INDIRECT`)  
errors TEXT : Errors encountered (stored as JSON or delimited text)  
ossf_score REAL : OpenSSF security score  

### dependency_edges
id INTEGER : Unique edge identifier (Primary Key)  
project_id TEXT : Associated project (Foreign Key to `project.id`)  
graph_id TEXT : Identifier for the dependency graph  
from_node_index INTEGER : Index of the source node  
to_node_index INTEGER : Index of the target node  
requirement TEXT : Dependency requirement (e.g., version constraints)  


# API

Supported Endpoints:

Add or Update Dependency
- POST /dependency/add
- Add or update a dependency for a project.

- Example request:
```json{
  "project_name": "example-project",
  "dependency": {
    "versionKey": {
      "system": "GO",
      "name": "example-dependency",
      "version": "v1.0.0"
    },
    "bundled": false,
    "relation": "DIRECT",
    "errors": []
  }
}
```
Get Specific Dependency
- GET /dependency/get/{projectName}/{dependencyName}
- Retrieve details of a specific dependency.
- Example response:
```json{
  "versionKey": {
    "system": "GO",
    "name": "example-dependency",
    "version": "v1.0.0"
  },
  "bundled": false,
  "relation": "DIRECT",
  "errors": [],
  "ossf_score": 8.5
}
```
Delete Dependency
- DELETE /dependency/delete/{projectName}/{dependencyName}
- Remove a dependency from a project.

List Dependencies
- GET /dependencies?name=example&min_score=7.0
- List dependencies with optional filters for name and score.
- Example response:
```json[
  {
    "versionKey": {
      "system": "GO",
      "name": "example-dependency",
      "version": "v1.0.0"
    },
    "bundled": false,
    "relation": "DIRECT",
    "errors": [],
    "ossf_score": 8.5
  }
]
```

Retrieve Project Dependencies
- GET /dependency/{projectName}
- Get a project's dependency graph.
- Example response:
```json{
  "project_name": "example-project",
  "dependencies": [
    {
      "id": "dependency-one",
      "score": 8.5
    }
  ]
}
```
