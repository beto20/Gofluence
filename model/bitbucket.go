package model

type Bitbucket struct {
	Branch                string
	Commit                string
	DeploymentEnvironment string
	RepoFullName          string
	Token                 string
	Url                   string
	ImageRemoteUrl        string // TODO: refactor this code because it doesn't make sense be here
}
