package models

type Version struct {
	VersionKey VersionKey `json:"versionKey"`
	IsDefault  bool       `json:"isDefault"` // Indicates if it's default version
}

type PackageKey struct {
	System string `json:"system"` // (e.g. GO)
	Name   string `json:"name"`
}

type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PackageVersions struct {
	PackageKey PackageKey `json:"packageKey"` // Unique identifier for the package
	Versions   []Version  `json:"versions"`   // All available versions
}
