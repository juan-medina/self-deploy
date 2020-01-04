package deploy

import (
	"log"
	"os"
)

func New() error {
	var err error = nil
	log.Println("starting deployment...")

	_, onK8s := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	if onK8s {
		err = newPipeline()
	} else {
		err = newJob()
	}

	log.Println("deployment completed")

	return err
}
