package deploy

import (
	"os"
)

const goPath = "/tmp"
const repo = "https://github.com/juan-medina/self-deploy.git"
const path = goPath + "/src/github.com/juan-medina/self-deploy"
const serviceName = "self-deploy"
const serviceVersion = "0.0.1"
const dockerFile = "Dockerfile"
const registryInternal = "localhost:32000"
const appPort = 5000

func newPipeline() error {
	var err error = nil

	_ = os.Setenv("GOPATH", goPath)

	err = clone(repo, path)

	if err == nil {
		err = test(path)
	}

	if err == nil {
		err = build(path)
	}

	if err == nil {
		err = dockerBuild(registryInternal, serviceName, serviceVersion, dockerFile, path)
	}

	if err == nil {
		err = dockerPush(registryInternal, serviceName)
	}

	if err == nil {
		err = k8sDeploy(registryInternal, serviceName, serviceVersion, appPort)
	}

	return err
}
