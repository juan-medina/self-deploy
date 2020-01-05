package deploy

type buildSettings struct {
	Name       string
	Version    string
	Path       string
	DockerFile string
}

type pipelineSettings struct {
	GoPath       string
	Repo         string
	Job          buildSettings
	App          buildSettings
	Registry     string
	InternalPort int
	ExternalPort int
}
