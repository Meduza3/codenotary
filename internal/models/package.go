package models

type Version struct {
	VersionKey VersionKey `json:"versionKey"`
	IsDefault  bool       `json:"isDefault"` 
}

type PackageKey struct {
	System string `json:"system"` 
	Name   string `json:"name"`
}

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackageVersions struct {
	PackageKey PackageKey `json:"packageKey"` 
	Versions   []Version  `json:"versions"`   
}
