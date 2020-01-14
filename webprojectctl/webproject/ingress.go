package main

import (
	"log"

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func createIngress(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-ing",
			Namespace: deploymentInput.Namespace,
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: deploymentInput.IngressDomainName,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: deploymentInput.DeploymentName + "-svc",
										ServicePort: intstr.FromInt(deploymentInput.PrimaryContainerPort),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, foundErr := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-ing", metav1.GetOptions{})
	if foundErr != nil {
		result, errIngress := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Create(ingress)
		if errIngress != nil {
			panic(errIngress)
		}
		log.Printf("Created Memcahed Deployment - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())

	}
}
