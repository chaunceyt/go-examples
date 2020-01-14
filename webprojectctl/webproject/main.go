package main

//@TODO
// - Check to see if object exists and update it. Similar to the kubectl apply -f filename.yaml
// - Sidecar support

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// WebProjectInput struct.
type WebProjectInput struct {
	DeploymentName           string
	PrimaryContainerName     string
	PrimaryContainerImageTag string
	PrimaryContainerPort     int
	Replicas                 int32
	Namespace                string
	CacheEngine              string
	DatabaseEngine           string
	DatabaseEngineImage      string
	IngressDomainName        string
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	deploymentName := flag.String("deployment-name", "", "Webproject Name")
	primaryContainerName := flag.String("primary-container-name", "", "Primary Container name")
	primaryContainerImageTag := flag.String("prinary-container-image-tag", "", "Primary Container image and tag")
	primaryContainerPort := flag.Int("primary-container-port", 8080, "Primary container port")
	replicas := flag.Int("replicas", 1, "Number of replicas")
	ingressDomainName := flag.String("domain-name", "", "Domainname for workload")
	projectNamespace := flag.String("namespace", "", "Project Namespace")
	flag.Parse()

	deploymentInput := WebProjectInput{
		DeploymentName:           *deploymentName,
		PrimaryContainerName:     *primaryContainerName,
		PrimaryContainerImageTag: *primaryContainerImageTag,
		PrimaryContainerPort:     *primaryContainerPort,
		Replicas:                 int32(*replicas),
		IngressDomainName:        *ingressDomainName,
		Namespace:                *projectNamespace,
	}

	//	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	//client, err := dynamic.NewForConfig(config)
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	createWebProject(client, deploymentInput)
}

func int32ptr(i int32) *int32 {
	return &i
}
