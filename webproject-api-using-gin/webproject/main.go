package main

//@TODO
// - Check to see if object exists and update it. Similar to the kubectl apply -f filename.yaml
// - Sidecar support

import (
	"flag"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var mySigningKey = []byte("captainjacksparrowsayshi")

type WebProjectInput struct {
	DeploymentName           string `json:"deploymentName"`
	PrimaryContainerName     string `json:"primaryContainerName"`
	PrimaryContainerImageTag string `json:"primaryContainerImageTag"`
	PrimaryContainerPort     int    `json:"primaryContainerPort"`
	Replicas                 int32  `json:"replicas"`
	Namespace                string `json:"namespace"`
	CacheEngine              string `json:"cacheEngine"`
	DatabaseEngine           string `json:"databaseEngine"`
	DatabaseEngineImage      string `json:"databaseEngineImage"`
	IngressDomainName        string `json:"ingressDomainName`
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	//client, err := dynamic.NewForConfig(config)
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.POST("/create-webproject", createWebProject(client))
	r.Run(":8084")
}

func int32ptr(i int32) *int32 {
	return &i
}
