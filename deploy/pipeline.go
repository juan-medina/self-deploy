package deploy

import "os"

const goPath = "/tmp"
const repo = "https://github.com/juan-medina/self-deploy.git"
const path = goPath + "/src/github.com/juan-medina/self-deploy"
const serviceName = "self-deploy"
const serviceVersion = "0.0.1"
const dockerFile = "Dockerfile"

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
		err = dockerBuild(registry, serviceName, serviceVersion, dockerFile, path)
	}

	return err
}
