package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func clone(url string) error {
	log.Println("cloning from git " + url + "...")
	out, err := exec.Command("git", "clone",
		url,
		"/tmp/self-deploy").CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error cloning software, err = , %s\n%s", err.Error(), out))
	}
	log.Println("cloning done")

	out2, _ := exec.Command("find", "/tmp").CombinedOutput()

	log.Println(fmt.Sprintf("output %q", out2))
	
	return err
}
