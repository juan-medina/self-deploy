package deploy

func newPipeline() error {
	var err error = nil

	err = clone("https://github.com/juan-medina/self-deploy.git")
	
	return err
}
