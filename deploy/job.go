package deploy

import (
	"log"
)

func newJob() error {
	log.Println("creating new job...")
	var err error = nil

	err = dockerBuild(settings.Registry, settings.Job.Name, settings.Job.Version, settings.Job.DockerFile, settings.Job.Path)

	if err == nil {
		err = dockerPush(settings.Registry, settings.Job.Name)
	}

	if err == nil {
		err = k8sJob(settings.Registry, settings.Job.Name, settings.Job.Version)
	}

	log.Println("job created")

	return err
}
