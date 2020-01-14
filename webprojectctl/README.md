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

## Example go run

```
 go run webproject/*.go \
  -deployment-name=vrt-manager-03 \
  -primary-container-name=vrt-manager \
  -prinary-container-image-tag=chaunceyt/vrt-manager-httpd \
  -primary-container-port=8080 \
  -replicas=1 \
  -domain-name=project-d.kube.domain.tld \
  -namespace=project-d
```



