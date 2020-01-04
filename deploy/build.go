package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func build(path string) error {
	log.Println("building software")
	cmd := exec.Command("go", "build")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error building software, err = , %s\n%s", err.Error(), string(out)))
	}
	log.Println("software built")

	return err
}
