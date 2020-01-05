package deploy

import (
	"errors"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"
	"log"
	"time"
)

func createNewDeployment(client v1.DeploymentInterface, service string, image string, port int32) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: service + "-dc",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": service,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": service,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  service + "-container",
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: port,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := client.Create(deployment)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating deployment = , %s\n", err.Error()))
	}

	return nil
}

func getClientSet() (*kubernetes.Clientset, error) {
	log.Println("getting k8s client..")

	var config, err = rest.InClusterConfig()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in kubernetes, err = , %s\n", err.Error()))
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error in kubernetes, err = , %s\n", err.Error()))
	}

	log.Println("got k8s client")

	return clientSet, nil
}

const interval = 10 * time.Second
const timeout = 10 * time.Minute

func waitDeletion(client v1.DeploymentInterface, deploymentName string) error {
	log.Println("wait for deletion")
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		exist, err := deploymentExist(client, deploymentName)
		log.Println(fmt.Sprintf("inside wait, exist =%v, err null =%v", exist, err == nil))
		if err != nil {
			return true, err
		}
		if exist {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return errors.New(fmt.Sprintf("error waiting on deployment deletion = , %s\n", err.Error()))
	}
	return nil
}

func deploymentExist(client v1.DeploymentInterface, deploymentName string) (bool, error) {
	_, err := client.Get(deploymentName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func deleteDeployment(client v1.DeploymentInterface, deploymentName string) error {
	log.Println("deleting deployment ...")

	deletePolicy := metav1.DeletePropagationForeground
	errDelete := client.Delete(deploymentName, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if errDelete != nil {
		return errors.New(fmt.Sprintf("error deleting deployment = , %s\n", errDelete.Error()))
	}

	log.Println("waiting for deletion ...")
	errWait := waitDeletion(client, deploymentName)
	if errWait != nil {
		return errors.New(fmt.Sprintf("error waiting deployment = , %s\n", errWait.Error()))
	}
	log.Println("waiting for deletion completed")

	log.Println("deployment deleted ...")

	return nil
}

func createDeployment(clientSet *kubernetes.Clientset, service string, image string, port int32) error {
	deploymentsClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	deploymentName := service + "-dc"
	exist, err := deploymentExist(deploymentsClient, deploymentName)
	if err != nil {
		return errors.New(fmt.Sprintf("error checking if deployment exist = , %s\n", err.Error()))
	}
	if exist {
		log.Println("deployment exists need to delete")
		errDelete := deleteDeployment(deploymentsClient, deploymentName)
		if errDelete != nil {
			return errors.New(fmt.Sprintf("error deleting deployment = , %s\n", errDelete.Error()))
		}
	} else {
		log.Println("deployment does not exist")
	}

	log.Println("creating deployment ...")
	errCreate := createNewDeployment(deploymentsClient, service, image, port)
	if errCreate != nil {
		return errors.New(fmt.Sprintf("error creating deployment = , %s\n", errCreate.Error()))
	}
	log.Println("deployment created")
	return nil
}

func k8sDeploy(registry string, serviceName string, version string, port int32) error {
	tag := registry + "/" + serviceName + ":" + version
	log.Println("deploying into k8s ...")
	clientSet, err := getClientSet()
	if err != nil {
		return errors.New(fmt.Sprintf("error getting client set, err = , %s\n", err.Error()))
	}

	err = createDeployment(clientSet, serviceName, tag, port)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating deployment, err = , %s\n", err.Error()))
	}

	log.Println("deployment into k8s completed")
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
