package types

type BuildSettings struct {
	Name       string
	Version    string
	Path       string
	DockerFile string
	Registry   string
}
type PipelineSettings struct {
	GoPath       string
	Repo         string
	BuildSettings
	InternalPort int
	ExternalPort int
}
