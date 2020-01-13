package main

import (
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
	_, errIngress := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Create(ingress)
	if errIngress != nil {
		panic(errIngress)
	}
}
