package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func dockerBuild(registry string, imgName string, version string, dockerFile string, path string) error {
	tag := registry + "/" + imgName + ":" + version
	log.Println("building docker " + tag + "...")
	cmd := exec.Command("docker", "build",
		"-t", tag,
		"-f", dockerFile,
		path)
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error building docker, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("docker " + tag + " built")
	return nil
}

func dockerPush(registry string, imgName string) error {
	fullImageName := registry + "/" + imgName
	log.Println("pushing docker image " + fullImageName + " ...")
	out, err := exec.Command("docker", "push", fullImageName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error pushing docker, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("docker for image " + fullImageName + " pushed")
	return nil
}
