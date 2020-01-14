package main

import (
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func createMemcachedWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	memcachedDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-memcached",
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app": deploymentInput.DeploymentName + "-memcached",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentInput.DeploymentName + "-memcached",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentInput.DeploymentName + "-memcached",
					Labels: map[string]string{
						"app": deploymentInput.DeploymentName + "-memcached",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "memcached",
							Image:           "memcached",
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								createContainerPort(6379),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-memcached", metav1.GetOptions{})
	if foundErr != nil {
		// Create  Memcached Deployment
		log.Println("Creating memcached deployment...")
		result, err := client.AppsV1().Deployments(deploymentInput.Namespace).Create(memcachedDeployment)
		if err != nil {
			panic(err)
		}
		// log.Printf("Created memcached deployment %q.\n", resultRedis.GetName())
		log.Printf("Created Memcahed Deployment - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())

	}

	serviceName := deploymentInput.DeploymentName + "-memcached-svc"
	labels := map[string]string{
		"app": deploymentInput.DeploymentName + "-memcached",
	}
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: v1.ServiceSpec{
			Selector: labels,
			Ports: []v1.ServicePort{{
				Port:       11211,
				TargetPort: intstr.FromInt(11211),
			}},
		},
	}

	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-memcached-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		log.Println("Creating memcached service...")
		service, errRedisService := client.CoreV1().Services(deploymentInput.Namespace).Create(service)
		if errRedisService != nil {
			panic(errRedisService)
		}
		log.Printf("Created Memcached Service - Name: %q, UID: %q\n", service.GetObjectMeta().GetName(), service.GetObjectMeta().GetUID())

	}
}
