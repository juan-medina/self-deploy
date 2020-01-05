package deploy

import (
	"github.com/jamedina/self-deploy/deploy/types"
	"os"
)

func newPipeline() error {
	settings := types.PipelineSettings{
		GoPath: "/tmp",
		Repo:   "https://github.com/juan-medina/self-deploy.git",
		BuildSettings: types.BuildSettings{
			Name:       "self-deploy",
			Version:    "0.0.1",
			Path:       "/tmp/src/github.com/juan-medina/self-deploy",
			DockerFile: "Dockerfile",
			Registry:   "localhost:32000",
		},
		InternalPort: 5000,
		ExternalPort: 8080,
	}

	var err error = nil

	err = os.Setenv("GOPATH", settings.GoPath)

	if err == nil {
		err = clone(settings.Repo, settings.Path)
	}

	if err == nil {
		err = test(settings.Path)
	}

	if err == nil {
		err = build(settings.Path)
	}

	if err == nil {
		err = dockerBuild(settings.BuildSettings)
	}

	if err == nil {
		err = dockerPush(settings.BuildSettings)
	}

	if err == nil {
		err = k8sDeploy(settings)
	}

	return err
}
