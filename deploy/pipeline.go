package deploy

import (
	"os"
)

const goPath = "/tmp"
const repo = "https://github.com/juan-medina/self-deploy.git"
const path = goPath + "/src/github.com/juan-medina/self-deploy"
const app = "self-deploy"
const version = "0.0.1"
const dockerFile = "Dockerfile"
const registry = "localhost:32000"
const internalPort = 5000
const externalPort = 8080

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
		err = dockerBuild(registry, app, version, dockerFile, path)
	}

	if err == nil {
		err = dockerPush(registry, app)
	}

	if err == nil {
		err = k8sDeploy(registry, app, version, internalPort, externalPort)
	}

	return err
}
