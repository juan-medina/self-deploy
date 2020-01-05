package deploy

import (
	"os"
)

var settings = pipelineSettings{
	GoPath: "/tmp",
	Repo:   "https://github.com/juan-medina/self-deploy.git",
	App: buildSettings{
		Name:       "self-deploy",
		Version:    "0.0.1",
		Path:       "/tmp/src/github.com/juan-medina/self-deploy",
		DockerFile: "Dockerfile",
	},
	Job: buildSettings{
		Name:       "self-deploy-job",
		Version:    "0.0.1",
		Path:       ".",
		DockerFile: "Dockerfile.job",
	},
	Registry:     "localhost:32000",
	InternalPort: 5000,
	ExternalPort: 8080,
}

func newPipeline() error {
	var err error = nil

	err = os.Setenv("GOPATH", settings.GoPath)

	if err == nil {
		err = clone(settings.Repo, settings.App.Path)
	}

	if err == nil {
		err = test(settings.App.Path)
	}

	if err == nil {
		err = build(settings.App.Path)
	}

	if err == nil {
		err = dockerBuild(settings.Registry, settings.App.Name, settings.App.Version, settings.App.DockerFile, settings.App.Path)
	}

	if err == nil {
		err = dockerPush(settings.Registry, settings.App.Name)
	}

	if err == nil {
		err = k8sDeploy(settings.Registry, settings.App.Name, settings.App.Version, settings.InternalPort, settings.ExternalPort)
	}

	return err
}
