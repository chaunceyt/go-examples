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

func createWebprojectWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// WebProject Deployment.
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName,
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app": deploymentInput.DeploymentName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(deploymentInput.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploymentInput.DeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentInput.DeploymentName,
					Labels: map[string]string{
						"app": deploymentInput.DeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentInput.PrimaryContainerName,
							Image:           deploymentInput.PrimaryContainerImageTag,
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts: []corev1.VolumeMount{
								createVolumeMount("webroot", "/var/www/webroot"),
								createVolumeMount("files", "/var/www/html/sites/default/files"),
							},
							Ports: []corev1.ContainerPort{
								createContainerPort(int32(deploymentInput.PrimaryContainerPort)),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						createEmptyDirVolume("webroot"),
						attachVolumeFromClaim("files", "webfiles", deploymentInput),
					},
				},
			},
		},
	}

	// Create Web Project Deployment
	foundWebProject, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName, metav1.GetOptions{})
	if foundErr != nil {
		log.Println("Creating webproject deployment...")
		resultWebProject, errWebProject := client.AppsV1().Deployments(deploymentInput.Namespace).Create(deployment)
		if errWebProject != nil {
			panic(errWebProject)
		}
		log.Printf("Created Deployment - Name: %q, UID: %q\n", resultWebProject.GetObjectMeta().GetName(), resultWebProject.GetObjectMeta().GetUID())
	} else {
		log.Println("Updating webproject deployment...")
		foundWebProject.Spec.Replicas = int32ptr(deploymentInput.Replicas)
		foundWebProject.Spec.Template.Spec.Containers[0].Image = deploymentInput.PrimaryContainerImageTag
		foundWebProjectResult, errFoundWebProject := client.AppsV1().Deployments(deploymentInput.Namespace).Update(foundWebProject)
		if errFoundWebProject != nil {
			panic(errFoundWebProject)
		}
		log.Printf("Updated Deployment - Name: %q, UID: %q\n", foundWebProjectResult.GetObjectMeta().GetName(), foundWebProjectResult.GetObjectMeta().GetUID())
	}

	log.Println("Creating service for WebProject.")
	serviceName := deploymentInput.DeploymentName + "-svc"

	webprojectLabels := map[string]string{
		"app": deploymentInput.DeploymentName,
	}
	webprojectService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: v1.ServiceSpec{
			Selector: webprojectLabels,
			Ports: []v1.ServicePort{{
				Port:       80,
				Protocol:   "TCP",
				TargetPort: intstr.FromInt(deploymentInput.PrimaryContainerPort),
			}},
		},
	}

	_, foundWebprojectServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-svc", metav1.GetOptions{})
	if foundWebprojectServiceErr != nil {
		webprojectServiceResp, errWebprojectService := client.CoreV1().Services(deploymentInput.Namespace).Create(webprojectService)
		if errWebprojectService != nil {
			panic(errWebprojectService)
		}
		log.Println(webprojectServiceResp.GetName())
		// log.Println("Created service for Webproject...")
		log.Printf("Created Webproject Service - Name: %q, UID: %q\n", webprojectServiceResp.GetObjectMeta().GetName(), webprojectServiceResp.GetObjectMeta().GetUID())

	}
}

func createContainerPort(portNumber int32) corev1.ContainerPort {
	return corev1.ContainerPort{
		ContainerPort: portNumber,
		Protocol:      corev1.ProtocolTCP,
	}
}
