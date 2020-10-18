# Go + helm3 sdk

Create kind cluster

kind create cluster --config kind-cluster.yaml

`kind-cluster.yaml`

```
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
  - |
    apiVersion: kubeadm.k8s.io/v1beta2
    kind: ClusterConfiguration
    metadata:
      name: config
    apiServer:
      extraArgs:
        "enable-admission-plugins": "NamespaceLifecycle,LimitRanger,ServiceAccount,TaintNodesByCondition,Priority,DefaultTolerationSeconds,DefaultStorageClass,PersistentVolumeClaimResize,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota"
nodes:
  - role: control-plane
  - role: worker
```

## Create helm chart

`helm create webapp`


`go run main.go -chart-path ./webapp -release webapp-v1`
