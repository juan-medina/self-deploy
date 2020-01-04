package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

const version = "0.0.1"
const jobImgName = "deploy-self-service"
const filesPath = "deploy/_files/job"

func jobDockerBuild() error {
	log.Println("building docker for job...")
	out, err := exec.Command("docker", "build", filesPath, "-t", "localhost:32000/"+jobImgName+":"+version).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error building docker, err = , %s\n%s", err.Error(), out))
	}
	log.Println("docker for job built")
	return nil
}

func jobDockerPush() error {
	log.Println("pushing docker for job...")
	out, err := exec.Command("docker", "push", "localhost:32000/"+jobImgName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error pushing docker, err = , %s\n%s", err.Error(), out))
	}
	log.Println("docker for job pushed")
	return nil
}

func createJob() error {

	log.Println("creating k8s job...")
	out, err := exec.Command("microk8s.kubectl", "create", "-f", filesPath+"/job.yml").CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error creating job, err = , %s\n%s", err.Error(), out))
	}
	log.Println("k8s job created")

	return nil
}

func NewJob() error {
	log.Println("creating new job...")
	var err error = nil

	err = jobDockerBuild()

	if err == nil {
		err = jobDockerPush()
	}

	if err == nil {
		err = createJob()
	}

	log.Println("job created")

	return err
}
