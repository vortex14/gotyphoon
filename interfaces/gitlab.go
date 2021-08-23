package interfaces


type GitlabProject struct {
	Name string `yaml:"name,omitempty"`
	Git string `yaml:"git,omitempty"`
	Id int	`yaml:"id,omitempty"`
}

type GitlabLabel struct {
	Id int `yaml:"id,omitempty"`
}

type GitlabClusterLabel struct {
	Url int `yaml:"url,omitempty"`
}

type GitLabel struct {
	Url string `yaml:"url,omitempty"`
	Remote string `yaml:"remote,omitempty"`
	Branch string `yaml:"branch,omitempty"`
}


type GitlabInterface interface {
	BuildCIResources()
}

type GitlabServer interface {
	GetAllProjectsList() []*GitlabProject
	SyncGitlabProjects()
	Deploy()
	HistoryPipelines()
}

