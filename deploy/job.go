package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

const version = "0.0.1"
const jobImgName = "deploy-self-service"
const jobDockerFile = "Dockerfile.job"
const registry = "localhost:32000"

func createJob() error {
	log.Println("creating k8s job...")
	out, err := exec.Command("microk8s.kubectl", "create", "-f", "job.yml").CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error creating job, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("k8s job created")

	return nil
}

func newJob() error {
	log.Println("creating new job...")
	var err error = nil

	err = dockerBuild(registry, jobImgName, version, jobDockerFile, ".")

	if err == nil {
		err = dockerPush(registry, jobImgName)
	}

	if err == nil {
		err = createJob()
	}

	log.Println("job created")

	return err
}
