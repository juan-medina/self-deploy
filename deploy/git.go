package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func clone(url string, path string) error {
	log.Println("cloning from git " + url + "...")
	out, err := exec.Command("git", "clone",
		url, path).CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error cloning software, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("cloning done")
	return err
}
