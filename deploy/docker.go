package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func dockerBuild(jobImgName string, version string, dockerFile string) error {
	log.Println("building docker for job...")
	out, err := exec.Command("docker", "build",
		"-t", "localhost:32000/"+jobImgName+":"+version,
		"-f", dockerFile,
		".").CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error building docker, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("docker for job built")
	return nil
}

func dockerPush(jobImgName string) error {
	log.Println("pushing docker for job...")
	out, err := exec.Command("docker", "push", "localhost:32000/"+jobImgName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error pushing docker, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("docker for job pushed")
	return nil
}
