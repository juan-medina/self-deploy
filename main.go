package main

import (
	"flag"
	"github.com/jamedina/self-deploy/deploy"
	"log"
	"net/http"
)

func main() {

	shouldDeploy := flag.Bool("deploy", false, "deploy this software")
	flag.Parse()

	if *shouldDeploy == false {
		server := NewServer()
		log.Println("starting server on port 5000")
		if err := http.ListenAndServe(":5000", server); err != nil {
			log.Fatalf("could not listen on port 5000 %v", err)
		}
	}else {
		err := deploy.New()
		if err!=nil {
			log.Fatalf("error deploying software %v", err)
		}
	}
}
