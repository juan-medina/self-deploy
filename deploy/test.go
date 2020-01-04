package deploy

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func test() error {
	log.Println("executing tests")
	out, err := exec.Command("go", "test", "/tmp/self-deploy").CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("error runing test, err = , %s\n%s", err.Error(), string(out)))
	}else{
		log.Println(string(out))
	}
	log.Println("test done")

	return err
}
