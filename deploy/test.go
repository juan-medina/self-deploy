package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func test(path string) error {
	log.Println("executing tests")
	cmd := exec.Command("go", "test", "-v")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error runing test, err = , %s\n%s", err.Error(), string(out)))
	}else{
		log.Println(string(out))
	}
	log.Println("test done")

	return err
}
