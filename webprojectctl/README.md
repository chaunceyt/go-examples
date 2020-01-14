# WebProject Control

As I continue to teach myself Golang. This hack session I decided to that [simple API endpoint](webproject-api-using-gin) and make it a cli tool.

Initial start was to move to flags. Later I plan to move to "github.com/urfave/cli" later.

## Issue

I work in an environment that uses GitLab + GKE as the CI solution for our developers. We're currently using a helm chart to generate the Kubernetes manifests for each prprojects workload. These workloads contain the following services

- Apache/PHP7.x PHP-FPM
- Mariadb
- Memcached
- Redis
- Solr
- ElasticSearch

## Proof of concept

Create a commandline tool that accepts parameters instructing it to create an environment with similar components stated above. This is ALPHA work.

## Example usage at the moment.

```
releases/webproject-ctl-darwin-amd64 \
  -deployment-name=vrt-manager \
  -primary-container-name=vrt-manager \
  -prinary-container-image-tag=chaunceyt/vrt-manager-httpd \
  -primary-container-port=8080 \
  -replicas=1 \
  -domain-name=project-d.kube.domain.tld \
  -namespace=webproject
```

## Output...

```
2020/01/13 20:40:37 Creating pvc...
2020/01/13 20:40:37 Created PVC - Name: "vrt-manager-webfiles-pvc", UID: "c49985c4-6857-4dd8-9bc0-1e37c39167d4"
2020/01/13 20:40:37 Unsupported CacheEngine selected or not defined
Creating webproject deployment...
Created Deployment - Name: "vrt-manager", UID: "6846f34c-717e-4c81-97d5-92b2e241a746"
Creating service for WebProject.
2020/01/13 20:40:37 Created Memcahed Deployment - Name: "vrt-manager-ing", UID: "4c5fccad-3fb1-4eac-8785-54212fefa1db"
Chaunceys-iMac:webprojectctl cthorn$ kubectl get po -n webproject
NAME                           READY   STATUS    RESTARTS   AGE
vrt-manager-6c9fdbdf69-7h2gs   1/1     Running   0          3s
Chaunceys-iMac:webprojectctl cthorn$ kubectl get all,pvc,ing -n webproject
NAME                               READY   STATUS    RESTARTS   AGE
pod/vrt-manager-6c9fdbdf69-7h2gs   1/1     Running   0          21s

NAME                      TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
service/vrt-manager-svc   ClusterIP   10.105.153.105   <none>        80/TCP    6m30s

NAME                          READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/vrt-manager   1/1     1            1           21s

NAME                                     DESIRED   CURRENT   READY   AGE
replicaset.apps/vrt-manager-6c9fdbdf69   1         1         1       21s

NAME                                             STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistentvolumeclaim/vrt-manager-webfiles-pvc   Bound    pvc-c49985c4-6857-4dd8-9bc0-1e37c39167d4   1Gi        RWO            standard       21s

NAME                                 HOSTS                       ADDRESS   PORTS   AGE
ingress.extensions/vrt-manager-ing   project-d.kube.domain.tld             80      21s
```



