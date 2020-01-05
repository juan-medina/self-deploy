package deploy

import (
	"log"
)

const jobApp = app + "-job"
const jobDockerFile = "Dockerfile.job"

func newJob() error {
	log.Println("creating new job...")
	var err error = nil

	err = dockerBuild(registry, jobApp, version, jobDockerFile, ".")

	if err == nil {
		err = dockerPush(registry, jobApp)
	}

	if err == nil {
		err = k8sJob(registry, jobApp, version)
	}

	log.Println("job created")

	return err
}
