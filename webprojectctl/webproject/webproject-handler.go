package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
)

// createWebProject - create all of the desire components.
func createWebProject(client *kubernetes.Clientset, deploymentInput WebProjectInput) {

	// Create Persistent Volume claims first.
	createPersistentVolumeClaim("webfiles", client, deploymentInput)

	// Determine if we are needing to deploy a database.
	var useDatabase bool

	if deploymentInput.DatabaseEngine == "" || deploymentInput.DatabaseEngineImage == "" {
		useDatabase = false
	} else {
		useDatabase = true
	}

	// Create database workload.
	if useDatabase == true {
		createPersistentVolumeClaim("db", client, deploymentInput)
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

}
