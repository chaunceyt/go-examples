package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/briandowns/spinner"

	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/storage/driver"
)

var helmConfig *action.Configuration
var chartRef *chart.Chart

var (
	chartPath string	
	namespace string
	releaseName string
)

func main() {
	fmt.Println("Helm client")

	flag.StringVar(&chartPath, "chart-path", "", "Path to helm chart")
    //flag.StringVar(&namespace, "namespace", "", "Namspace")
	flag.StringVar(&releaseName, "release", "", "release name")
	flag.Parse()
	namespace = "webapp2"
	s := spinner.New(spinner.CharSets[14], 125*time.Millisecond)
	s.Suffix = "Creating project..."
	s.Start()
	

	settings := cli.New()
	helmDriver := "secret"

	actionConfig := new(action.Configuration)

	
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)	
	}
	helmConfig = new(action.Configuration)
	if err := helmConfig.Init(settings.RESTClientGetter(), namespace, helmDriver, func(format string, v ...interface{}) {}); err != nil {
		os.Exit(1)
	}
    // define values
    vals := map[string]interface{}{
        "image": map[string]interface{}{
            "tag": "1.17.1",
            },
    }
   fmt.Println(string(*helmConfig))
	// load chart from the path 
	chart, err := loader.Load(chartPath)
	if err != nil {
		panic(err)
	}

	// Check to see if the release exists in the namespace.
	statusCommand := action.NewStatus(helmConfig)
	status, err := statusCommand.Run(releaseName)
		// if there's a non-404 error, something went wrong and we'll exit out
	if err != nil && err != driver.ErrReleaseNotFound {
		fmt.Println("Failed to retrieve Helm release", err.Error())
		os.Exit(1)
	}

	if err == nil && status != nil {

		fmt.Printf("Release %q exists. Upgrading in namespace %q.\n", releaseName, namespace)
		client := action.NewUpgrade(actionConfig)
		client.Namespace = namespace
		// client.DryRun = true - very handy!
	
		// install the chart here
		rel, err := client.Run(releaseName, chart, vals)
		if err != nil {
			panic(err)
		}
	
		log.Printf("Chart Upgraded from path: %s in namespace: %s\n", rel.Name, rel.Namespace)
		// this will confirm the values set during installation
		log.Println(rel.Config)

	} else {

		client := action.NewInstall(actionConfig)
		client.Namespace = namespace
		client.ReleaseName = releaseName
		// client.DryRun = true - very handy!
	
		// install the chart here
		rel, err := client.Run(chart, vals)
		if err != nil {
			panic(err)
		}
	
		log.Printf("Installed Chart from path: %s in namespace: %s\n", rel.Name, rel.Namespace)
		// this will confirm the values set during installation
		log.Println(rel.Config)
	}
	
	s.Stop()

}