package models

type Node struct {
	VersionKey VersionKey `json:"versionKey"`
	Bundled    bool       `json:"bundled"`
	Relation   string     `json:"relation"`
	Errors     []string   `json:"errors"`
}

type Edge struct {
	FromNode    int    `json:"fromNode"`
	ToNode      int    `json:"toNode"`
	Requirement string `json:"requirement"`
}

type DependencyGraph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
	Error string `json:"error"`
}
