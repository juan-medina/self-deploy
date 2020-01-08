package deploy

import (
	"errors"
	"fmt"
	"github.com/jamedina/self-deploy/deploy/types"
	appsv1 "k8s.io/api/apps/v1"
	v13 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
	v12 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"time"
)

func createNewDeployment(client v1.DeploymentInterface, app string, image string, port int) error {
	replicas := int32(1)
	containerPort := int32(port)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: app,
			Labels: map[string]string{
				"app": app,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            app + "-container",
							Image:           image,
							ImagePullPolicy: apiv1.PullAlways,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: containerPort,
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

func getClientSetOnCluster() (*kubernetes.Clientset, error) {
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

func waitDeploymentDeletion(client v1.DeploymentInterface, app string) error {
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		exist, err := deploymentExist(client, app)
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

func deploymentExist(client v1.DeploymentInterface, app string) (bool, error) {
	_, err := client.Get(app, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func deleteDeployment(client v1.DeploymentInterface, app string) error {
	log.Println("deleting deployment ...")

	deletePolicy := metav1.DeletePropagationForeground
	errDelete := client.Delete(app, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if errDelete != nil {
		return errors.New(fmt.Sprintf("error deleting deployment = , %s\n", errDelete.Error()))
	}

	log.Println("waiting for deletion ...")
	errWait := waitDeploymentDeletion(client, app)
	if errWait != nil {
		return errors.New(fmt.Sprintf("error waiting deployment = , %s\n", errWait.Error()))
	}
	log.Println("waiting for deletion completed")

	log.Println("deployment deleted ...")

	return nil
}

func createDeployment(clientSet *kubernetes.Clientset, app string, image string, port int) error {
	deploymentsClient := clientSet.AppsV1().Deployments(apiv1.NamespaceDefault)
	exist, err := deploymentExist(deploymentsClient, app)
	if err != nil {
		return errors.New(fmt.Sprintf("error checking if deployment exist = , %s\n", err.Error()))
	}
	if exist {
		log.Println("deployment exists need to delete")
		errDelete := deleteDeployment(deploymentsClient, app)
		if errDelete != nil {
			return errors.New(fmt.Sprintf("error deleting deployment = , %s\n", errDelete.Error()))
		}
	} else {
		log.Println("deployment does not exist")
	}

	log.Println("creating deployment ...")
	errCreate := createNewDeployment(deploymentsClient, app, image, port)
	if errCreate != nil {
		return errors.New(fmt.Sprintf("error creating deployment = , %s\n", errCreate.Error()))
	}
	log.Println("deployment created")
	return nil
}

func waitServiceDeletion(client v12.ServiceInterface, app string) error {
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		exist, err := serviceExist(client, app)
		if err != nil {
			return true, err
		}
		if exist {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return errors.New(fmt.Sprintf("error waiting on service deletion = , %s\n", err.Error()))
	}
	return nil
}

func serviceExist(client v12.ServiceInterface, app string) (bool, error) {
	_, err := client.Get(app, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func createNewService(client v12.ServiceInterface, app string, port int, target int) error {
	serviceSpec := &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: app,
			Labels: map[string]string{
				"app": app,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       app,
					Port:       int32(port),
					TargetPort: intstr.FromInt(target),
					Protocol:   apiv1.ProtocolTCP,
				},
			},
			Selector:        map[string]string{"app": app},
			Type:            apiv1.ServiceTypeClusterIP,
			SessionAffinity: apiv1.ServiceAffinityNone,
		},
	}
	_, err := client.Create(serviceSpec)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating service = , %s\n", err.Error()))
	}
	return nil
}

func deleteService(client v12.ServiceInterface, app string) error {
	log.Println("deleting service ...")

	deletePolicy := metav1.DeletePropagationForeground
	errDelete := client.Delete(app, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if errDelete != nil {
		return errors.New(fmt.Sprintf("error deleting service = , %s\n", errDelete.Error()))
	}

	log.Println("waiting for deletion ...")
	errWait := waitServiceDeletion(client, app)
	if errWait != nil {
		return errors.New(fmt.Sprintf("error waiting service = , %s\n", errWait.Error()))
	}
	log.Println("waiting for deletion completed")

	log.Println("service deleted ...")

	return nil
}

func createService(clientSet *kubernetes.Clientset, app string, port int, target int) error {
	serviceClient := clientSet.CoreV1().Services(apiv1.NamespaceDefault)

	exist, err := serviceExist(serviceClient, app)
	if err != nil {
		return errors.New(fmt.Sprintf("error checking if service exist = , %s\n", err.Error()))
	}
	if exist {
		log.Println("service exists need to delete")
		errDelete := deleteService(serviceClient, app)
		if errDelete != nil {
			return errors.New(fmt.Sprintf("error deleting service = , %s\n", errDelete.Error()))
		}
	} else {
		log.Println("service does not exist")
	}

	log.Println("creating service ...")
	errCreate := createNewService(serviceClient, app, port, target)
	if errCreate != nil {
		return errors.New(fmt.Sprintf("error creating service = , %s\n", errCreate.Error()))
	}
	log.Println("service created")
	return nil
}

func k8sDeploy(appSettings types.PipelineSettings) error {
	tag := appSettings.Registry + "/" + appSettings.Name + ":" + appSettings.Version
	log.Println("deploying into k8s ...")
	clientSet, err := getClientSetOnCluster()
	if err != nil {
		return errors.New(fmt.Sprintf("error getting client set, err = , %s\n", err.Error()))
	}

	err = createDeployment(clientSet, appSettings.Name, tag, appSettings.InternalPort)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating deployment, err = , %s\n", err.Error()))
	}

	log.Println("deployment into k8s completed")

	err = createService(clientSet, appSettings.Name, appSettings.ExternalPort, appSettings.InternalPort)
	if err != nil {
		return errors.New(fmt.Sprintf("error creating service, err = , %s\n", err.Error()))
	}

	return nil
}

func getClientSetOnLocal() (*kubernetes.Clientset, error) {
	log.Println("getting k8s client..")

	home := os.Getenv("HOME")
	kubeConfig := filepath.Join(home, ".kube", "config")

	var config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
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

func createJob(clientSet *kubernetes.Clientset, registry string, app string, version string) error {
	tag := registry + "/" + app + ":" + version
	jobsClient := clientSet.BatchV1().Jobs(apiv1.NamespaceDefault)
	privileged := true
	file := apiv1.HostPathFile

	job := &v13.Job{
		TypeMeta: metav1.TypeMeta{
			Kind: "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: app + "-job",
			Labels: map[string]string{
				"app":       app,
				"job-group": app,
			},
		},
		Spec: v13.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       app,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            app + "job-container",
							Image:           tag,
							ImagePullPolicy: apiv1.PullAlways,
							SecurityContext: &apiv1.SecurityContext{
								Privileged: &privileged,
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "docker-socket-volume",
									MountPath: "/var/run/docker.sock",
								},
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyNever,
					Volumes: []apiv1.Volume{
						{
							Name: "docker-socket-volume",
							VolumeSource: apiv1.VolumeSource{
								HostPath: &apiv1.HostPathVolumeSource{
									Path: "/var/run/docker.sock",
									Type: &file,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := jobsClient.Create(job)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating job, err = , %s\n", err.Error()))
	}

	return nil
}

func k8sJob(buildSettings types.BuildSettings) error {
	clientSet, err := getClientSetOnLocal()

	if err != nil {
		return errors.New(fmt.Sprintf("error getting client set, err = , %s\n", err.Error()))
	}

	err = createJob(clientSet, buildSettings.Registry, buildSettings.Name, buildSettings.Version)

	if err != nil {
		return errors.New(fmt.Sprintf("error creating job, err = , %s\n", err.Error()))
	}

	return nil
}
