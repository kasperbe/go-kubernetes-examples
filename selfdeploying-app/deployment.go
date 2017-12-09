package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
)

var deploymentconfig = "Kube/deployment.yaml"

// Deployment - configuration for final deployment from yaml file
func Deployment(replicas int32) *appsv1beta1.Deployment {

	// Find absolute path of binary
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Read YAML and convert to JSON
	kubedata, err := ioutil.ReadFile(dir + "/" + deploymentconfig)
	if err != nil {
		log.Fatal(err)
	}
	kubedata, err = yaml.YAMLToJSON(kubedata)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal JSON to struct
	dep := appsv1beta1.Deployment{}
	json.Unmarshal(kubedata, &dep)

	return &dep
}
