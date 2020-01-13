package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

// createWebProject - create all of the desire components.
func createWebProject(client *kubernetes.Clientset) gin.HandlerFunc {
	return func(c *gin.Context) {
		deploymentInput := WebProjectInput{}
		c.Bind(&deploymentInput)

		// Create Persistent Volume claims first.
		createPersistentVolumeClaim("webfiles", client, deploymentInput)
		createPersistentVolumeClaim("db", client, deploymentInput)

		// Determine if we are needing to deploy a database.
		var useDatabase bool

		if deploymentInput.DatabaseEngine == "" || deploymentInput.DatabaseEngineImage == "" {
			useDatabase = false
		} else {
			useDatabase = true
		}

		// Create database workload.
		if useDatabase == true {
			createDatabaseWorkload(client, deploymentInput)
		}

		// Select the cacheEngine.
		if deploymentInput.CacheEngine == "redis" {
			// Using Redis for CacheEngine
			createRedisWorkload(client, deploymentInput)

		} else if deploymentInput.CacheEngine == "memcached" {
			createMemcachedWorkload(client, deploymentInput)

		} else {
			log.Println("Unsupported CacheEngine selected or not defined")
		}

		// Create project's primary workload.
		createWebprojectWorkload(client, deploymentInput)

		// Setup domain(s) for Webproject
		createIngress(client, deploymentInput)

		c.JSON(http.StatusOK, "success")
	}

}
