# Self Deploy

Software that is able to deploy itself

[![License: Apache2](https://img.shields.io/badge/license-Apache%202-blue.svg)](/LICENSE)

## info

This is and experiment about creating software that is able to deploy it self, to understand further this concept, and the code that it use you could review this presentation : 

https://juan-medina.github.io/self-deploy/

## building the software

You could just follow the bellow instructions
```shell script
$ git clone https://github.com/juan-medina/self-deploy.git
$ cd self-deploy
$ go build
```
## running the example

To just run the example locally we could just do : 

```shell script
$ ./self-deploy 
2020/01/06 17:07:40 starting server on port 5000

```
## testing it locally

To test our service we could use [HTTPie](https://httpie.org/) :

```shell script
$ http :5000/random
HTTP/1.1 200 OK
Content-Length: 19
Content-Type: text/plain; charset=utf-8
Date: Mon, 06 Jan 2020 17:15:53 GMT

5577006791947779410

$ http :5000/random
HTTP/1.1 200 OK
Content-Length: 19
Content-Type: text/plain; charset=utf-8
Date: Mon, 06 Jan 2020 17:15:55 GMT

8674665223082153551
```
## deploying into kubernetes

For this we will need : 

- A kubernetes cluster, such MicroK8s or minikube
- Cluster configured on your machine, aka `.kube/config`
- Docker-ce installed
- A docker registry available at `localhost:32000`
- Kubernetes host configured to speak with docker daemon and with the registry

Now we could launch the deployment with

```shell script
$ ./self-deploy -deploy
2020/01/06 00:19:37 starting deployment...
2020/01/06 00:19:37 creating new job...
2020/01/06 00:19:37 building docker localhost:32000/self-deploy-job:0.0.1...
2020/01/06 00:19:39 docker localhost:32000/self-deploy-job:0.0.1 built
2020/01/06 00:19:39 pushing docker image localhost:32000/self-deploy-job ...
2020/01/06 00:19:39 docker for image localhost:32000/self-deploy-job pushed
2020/01/06 00:19:39 getting k8s client..
2020/01/06 00:19:39 got k8s client
2020/01/06 00:19:39 job created
2020/01/06 00:19:39 deployment completed
```

## checking our cluster

After the deployment we could inspect what we have in our cluster

```shell script
$ microk8s.kubectl get all                            
NAME                                 READY   STATUS      RESTARTS   AGE
pod/self-deploy-86cb47fc7-h6x4m      1/1     Running     0          113s
pod/self-deploy-job-jobz5t6b-cj7hf   0/1     Completed   0          2m15s

NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/kubernetes    ClusterIP   10.152.183.1     <none>        443/TCP    21h
service/self-deploy   ClusterIP   10.152.183.222   <none>        8080/TCP   113s

NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/self-deploy   1/1     1            1           113s

NAME                                    DESIRED   CURRENT   READY   AGE
replicaset.apps/self-deploy-86cb47fc7   1         1         1       113s

NAME                                 COMPLETIONS   DURATION   AGE
job.batch/self-deploy-job-jobz5t6b   1/1           22s        2m15s

```

This application has first create a k8s.job that contains it own executable.

The job has run a self-contain pipeline that has deploy itself creating a k8s deployment and a service.

## checking the job log

We could use kubectl to check the job log

```shell script
$ kubectl logs job.batch/self-deploy-job-jobz5t6b
2020/01/06 00:19:40 starting deployment...
2020/01/06 00:19:40 cloning from git https://github.com/juan-medina/self-deploy.git...
2020/01/06 00:19:41 cloning done
2020/01/06 00:19:41 executing tests
2020/01/06 00:19:57 go: downloading k8s.io/apimachinery v0.17.0
go: downloading k8s.io/api v0.17.0
go: downloading k8s.io/client-go v0.17.0
go: extracting k8s.io/apimachinery v0.17.0
go: downloading github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
go: downloading k8s.io/klog v1.0.0
go: extracting k8s.io/klog v1.0.0
go: extracting k8s.io/client-go v0.17.0
go: downloading github.com/spf13/pflag v1.0.5
go: extracting github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
go: extracting github.com/spf13/pflag v1.0.5
go: extracting k8s.io/api v0.17.0
go: downloading gopkg.in/inf.v0 v0.9.1
go: downloading github.com/imdario/mergo v0.3.5
go: downloading golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
go: downloading github.com/golang/protobuf v1.3.2
go: downloading github.com/google/gofuzz v1.0.0
go: downloading github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
go: downloading github.com/davecgh/go-spew v1.1.1
go: downloading golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
go: downloading golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
go: downloading k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
go: extracting golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
go: downloading golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
go: extracting github.com/davecgh/go-spew v1.1.1
go: downloading github.com/json-iterator/go v1.1.8
go: extracting golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
go: extracting github.com/imdario/mergo v0.3.5
go: extracting gopkg.in/inf.v0 v0.9.1
go: extracting github.com/google/gofuzz v1.0.0
go: downloading github.com/modern-go/reflect2 v1.0.1
go: downloading sigs.k8s.io/yaml v1.1.0
go: extracting github.com/golang/protobuf v1.3.2
go: extracting k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
go: extracting github.com/modern-go/reflect2 v1.0.1
go: downloading github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
go: extracting sigs.k8s.io/yaml v1.1.0
go: downloading gopkg.in/yaml.v2 v2.2.4
go: extracting github.com/json-iterator/go v1.1.8
go: extracting gopkg.in/yaml.v2 v2.2.4
go: extracting github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
go: extracting golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
go: downloading golang.org/x/text v0.3.2
go: extracting golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
go: downloading golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456
go: extracting github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
go: extracting golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456
go: extracting golang.org/x/text v0.3.2
go: finding k8s.io/apimachinery v0.17.0
go: finding k8s.io/api v0.17.0
go: finding k8s.io/client-go v0.17.0
go: finding github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
go: finding github.com/google/gofuzz v1.0.0
go: finding github.com/imdario/mergo v0.3.5
go: finding golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
go: finding k8s.io/klog v1.0.0
go: finding github.com/spf13/pflag v1.0.5
go: finding github.com/golang/protobuf v1.3.2
go: finding golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
go: finding github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
go: finding gopkg.in/inf.v0 v0.9.1
go: finding golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456
go: finding github.com/json-iterator/go v1.1.8
go: finding github.com/modern-go/reflect2 v1.0.1
go: finding golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
go: finding sigs.k8s.io/yaml v1.1.0
go: finding github.com/davecgh/go-spew v1.1.1
go: finding k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
go: finding golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
go: finding golang.org/x/text v0.3.2
go: finding github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd
go: finding gopkg.in/yaml.v2 v2.2.4
=== RUN   TestRandom
--- PASS: TestRandom (0.00s)
=== RUN   TestNotFound
--- PASS: TestNotFound (0.00s)
PASS
ok  	github.com/jamedina/self-deploy	0.004s

2020/01/06 00:19:57 test done
2020/01/06 00:19:57 building software
2020/01/06 00:19:58 software built
2020/01/06 00:19:58 building docker localhost:32000/self-deploy:0.0.1...
2020/01/06 00:19:59 docker localhost:32000/self-deploy:0.0.1 built
2020/01/06 00:19:59 pushing docker image localhost:32000/self-deploy ...
2020/01/06 00:20:01 docker for image localhost:32000/self-deploy pushed
2020/01/06 00:20:01 deploying into k8s ...
2020/01/06 00:20:01 getting k8s client..
2020/01/06 00:20:01 got k8s client
2020/01/06 00:20:01 deployment does not exist
2020/01/06 00:20:01 creating deployment ...
2020/01/06 00:20:01 deployment created
2020/01/06 00:20:01 deployment into k8s completed
2020/01/06 00:20:01 service does not exist
2020/01/06 00:20:01 creating service ...
2020/01/06 00:20:01 service created
2020/01/06 00:20:01 deployment completed
```

Since this job is running on the cluster we could use fluent / kibana to view the logs of our jobs

# testing the service on the cluster

First we need the service ip address

```shell script
$ kubectl get svc/self-deploy -o jsonpath='{.spec.clusterIP}'
10.152.183.222
```

The we could just use [HTTPie](https://httpie.org/)

```shell script
$ http 10.152.183.222:8080/random
HTTP/1.1 200 OK
Content-Length: 19
Content-Type: text/plain; charset=utf-8
Date: Mon, 06 Jan 2020 00:28:48 GMT

5577006791947779410

$ http 10.152.183.222:8080/random
HTTP/1.1 200 OK
Content-Length: 19
Content-Type: text/plain; charset=utf-8
Date: Mon, 06 Jan 2020 00:28:49 GMT

8674665223082153551
```