package model

type Document struct {
	Name         string
	Version      string
	Parent       Parent
	Dependencies Dependencies
}

type Dependency struct {
	GroupId    string
	ArtifactId string
	Version    string
	Scope      string
}

type Dependencies struct {
	Dependency []Dependency
}

type Parent struct {
	Name    string
	Version string
}
