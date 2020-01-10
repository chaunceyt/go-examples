package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type WebProjectInput struct {
	DeploymentName           string `json:"deploymentName"`
	PrimaryContainerName     string `json:"primaryContainerName"`
	PrimaryContainerImageTag string `json:"primaryContainerImageTag"`
	PrimaryContainerPort     int    `json:"primaryContainerPort"`
	Replicas                 int32  `json:"replicas"`
	Namespace                string `json:"namespace"`
	CacheEngine              string `json:"cacheEngine"`
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

func createWebProject(client *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		deploymentInput := WebProjectInput{}
		c.Bind(&deploymentInput)

		//servicesClient := client.CoreV1().Services(corev1.NamespaceDefault)
		createPVC(client, deploymentInput)

		//ingress := *v1beta1.ExtensionsV1beta1().Ingress

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
									{
										Name:      "files",
										MountPath: "/var/www/html/sites/default/files",
									},
									{
										Name:      "webroot",
										MountPath: "/var/www/webroot",
									},
								},
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 80,
										Protocol:      corev1.ProtocolTCP,
									},
								},
							},
						},
						RestartPolicy: corev1.RestartPolicyAlways,
						Volumes: []corev1.Volume{
							GetSiteFilesVolume(deploymentInput),
							GetWebRootVolume(),
						},
					},
				},
			},
		}

		if deploymentInput.CacheEngine == "redis" {
			// Using Redis for CacheEngine
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
										{
											ContainerPort: 6379,
											Protocol:      corev1.ProtocolTCP,
										},
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
			resultRedis, errRedis := client.AppsV1().Deployments(deploymentInput.Namespace).Create(redisDeployment)
			if errRedis != nil {
				panic(errRedis)
			}
			log.Printf("Created redis deployment %q.\n", resultRedis.GetName())

			// move to a single func.
			serviceName := deploymentInput.DeploymentName + "-redis-svc"
			labels := map[string]string{
				"app": deploymentInput.DeploymentName + "-redis",
			}
			service := &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: serviceName,
				},
				Spec: v1.ServiceSpec{
					Selector: labels,
					Ports: []v1.ServicePort{{
						Port:       6379,
						TargetPort: intstr.FromInt(6379),
					}},
				},
			}
			service, errRedisService := client.CoreV1().Services(deploymentInput.Namespace).Create(service)
			if errRedisService != nil {
				panic(errRedisService)
			}
			//log.Println(service)

		} else if deploymentInput.CacheEngine == "memcached" {

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
										{
											ContainerPort: 6379,
											Protocol:      corev1.ProtocolTCP,
										},
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyAlways,
						},
					},
				},
			}
			// Create  Memcached Deployment
			log.Println("Creating memcached deployment...")
			resultRedis, errRedis := client.AppsV1().Deployments(deploymentInput.Namespace).Create(memcachedDeployment)
			if errRedis != nil {
				panic(errRedis)
			}
			log.Printf("Created memcached deployment %q.\n", resultRedis.GetName())
		} else {
			log.Println("Unsupported CacheEngine selected or not defined")
		}

		// Create Web Project Deployment
		log.Println("Creating deployment...")
		resultWebProject, errWebProject := client.AppsV1().Deployments(deploymentInput.Namespace).Create(deployment)
		// resultWebProject, errWebProject := client.Resource(deploymentRes).Namespace(deploymentInput.Namespace).Create(deployment, metav1.CreateOptions{})
		if errWebProject != nil {
			panic(errWebProject)
		}
		log.Printf("Created deployment %q.\n", resultWebProject.GetName())

		log.Println("Creating service for WebProject.")
		serviceName := deploymentInput.DeploymentName + "-svc"
		labels := map[string]string{
			"app": deploymentInput.DeploymentName,
		}
		webprojectService := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceName,
			},
			Spec: v1.ServiceSpec{
				Selector: labels,
				Ports: []v1.ServicePort{{
					Port:       80,
					TargetPort: intstr.FromInt(80),
				}},
			},
		}
		webprojectServiceResp, errWebprojectService := client.CoreV1().Services(deploymentInput.Namespace).Create(webprojectService)
		if errWebprojectService != nil {
			panic(errWebprojectService)
		}
		log.Println(webprojectServiceResp.GetName())
		log.Println("Created service for Webproject...")

		ingress := &v1beta1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      deploymentInput.DeploymentName + "-ing",
				Namespace: deploymentInput.Namespace,
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "domain.tld",
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

		c.JSON(http.StatusOK, "success")
		// c.Status(http.StatusNoContent)
	}

}

func int32ptr(i int32) *int32 {
	return &i
}

func createPVC(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentInput.DeploymentName + "-pvc",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("3Gi"),
				},
			},
		},
	}
	log.Println("Creating pvc...")
	resultPVC, errPVC := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Create(pvc)
	if errPVC != nil {
		panic(errPVC)
	}
	log.Printf("Created pvc %q.\n", resultPVC.GetName())

}

func GetWebRootVolume() corev1.Volume {
	return corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
		Name: "webroot",
	}
}

func GetSiteFilesVolume(deploymentInput WebProjectInput) corev1.Volume {
	return corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: deploymentInput.DeploymentName + "-pvc",
			},
		},
		Name: "files",
	}
}
