package deploy

import (
	"github.com/jamedina/self-deploy/deploy/types"
	"log"
)

func newJob() error {
	settings := types.BuildSettings{
		Name:       "self-deploy-job",
		Version:    "0.0.1",
		Path:       ".",
		DockerFile: "Dockerfile.job",
		Registry:   "localhost:32000",
	}

	log.Println("creating new job...")
	var err error = nil

	err = dockerBuild(settings)

	if err == nil {
		err = dockerPush(settings)
	}

	if err == nil {
		err = k8sJob(settings)
	}

	log.Println("job created")
	return err
}
