# WebProject API

As I continue to teach myself Golang. This hack session I decided to try [gin](https://github.com/gin-gonic/gin)  a HTTP web framework to create a simple API endpoint.

## Issue

I work in an environment that uses GitLab + GKE as the CI solution for our developers. We're currently using a helm chart to generate the Kubernetes manifests for each prprojects workload. These workloads contain the following services

- Apache/PHP7.x PHP-FPM
- Mariadb
- Memcached
- Redis
- Solr
- ElasticSearch

## Proof of concept

Create an API endpoint that accepts a JSON payload and creates an environment

## Example payload

```
{
	"deploymentName": "vrt-manager-httpd",
	"primaryContainerName": "vrt-manager",
	"primaryContainerImageTag": "chaunceyt/vrt-manager-httpd",
	"primaryContainerPort": 8080,
	"replicas": 1,
	"namespace": "default",
	"cacheEngine": "redis"
}
```



