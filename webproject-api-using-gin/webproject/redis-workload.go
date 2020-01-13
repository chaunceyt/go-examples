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

func createRedisWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	redisDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-redis",
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app": deploymentInput.DeploymentName + "-redis",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentInput.DeploymentName + "-redis",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentInput.DeploymentName + "-redis",
					Labels: map[string]string{
						"app": deploymentInput.DeploymentName + "-redis",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "redis",
							Image:           "redis",
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

	// Create  Redis Deployment
	log.Println("Creating redis deployment...")
	result, err := client.AppsV1().Deployments(deploymentInput.Namespace).Create(redisDeployment)
	if err != nil {
		panic(err)
	}
	log.Printf("Created Redis Deployment - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())

	serviceName := deploymentInput.DeploymentName + "-redis-svc"
	redisLabels := map[string]string{
		"app": deploymentInput.DeploymentName + "-redis",
	}
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: v1.ServiceSpec{
			Selector: redisLabels,
			Ports: []v1.ServicePort{{
				Port:       6379,
				TargetPort: intstr.FromInt(6379),
			}},
		},
	}

	log.Println("Creating redis service...")
	service, errRedisService := client.CoreV1().Services(deploymentInput.Namespace).Create(service)
	if errRedisService != nil {
		panic(errRedisService)
	}
	log.Printf("Created Redis Service - Name: %q, UID: %q\n", service.GetObjectMeta().GetName(), service.GetObjectMeta().GetUID())

}
