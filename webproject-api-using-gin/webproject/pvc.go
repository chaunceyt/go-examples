package main

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// createPersistentVolumeClaim - creates persistent volume claim
func createPersistentVolumeClaim(pvcType string, client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentInput.DeploymentName + "-" + pvcType + "-pvc",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	// Check to see if persistent volume claim exists already.
	_, foundErr := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-"+pvcType+"-pvc", metav1.GetOptions{})
	if foundErr != nil {
		log.Println("Creating pvc...")
		result, err := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Create(pvc)
		if err != nil {
			panic(err)
		}
		log.Printf("Created PVC - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())
	}
}

// createEmptyDirVolume
func createEmptyDirVolume(name string) corev1.Volume {
	return corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
		Name: name,
	}
}

// createVolumeMount - should be used to attach a volume to a deplopyment.
func createVolumeMount(name string, mountPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		MountPath: mountPath,
	}
}

// attachVolumeClaim attaches the Persistent volume to container.
func attachVolumeFromClaim(name string, pvType string, deploymentInput WebProjectInput) corev1.Volume {
	return corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: deploymentInput.DeploymentName + "-" + pvType + "-pvc",
			},
		},
		Name: name,
	}
}
