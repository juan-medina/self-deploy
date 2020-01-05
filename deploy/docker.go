package deploy

import (
	"errors"
	"fmt"
	"github.com/jamedina/self-deploy/deploy/types"
	"log"
	"os/exec"
)

func dockerBuild(settings types.BuildSettings) error {
	tag := settings.Registry + "/" + settings.Name + ":" + settings.Version
	log.Println("building docker " + tag + "...")
	cmd := exec.Command("docker", "build",
		"-t", tag,
		"-f", settings.DockerFile,
		settings.Path)
	cmd.Dir = settings.Path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error building docker, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("docker " + tag + " built")
	return nil
}

func dockerPush(settings types.BuildSettings) error {
	fullImageName := settings.Registry + "/" + settings.Name
	log.Println("pushing docker image " + fullImageName + " ...")
	out, err := exec.Command("docker", "push", fullImageName).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error pushing docker, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("docker for image " + fullImageName + " pushed")
	return nil
}
