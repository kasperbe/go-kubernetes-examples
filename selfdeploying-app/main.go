package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/pkg/errors"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func deployKubernetes(configPath string, replicas int32) error {
	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return errors.Wrap(err, "Couldn't build configuration")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "Couldn't connect to kubernetes")
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
	d := Deployment(replicas)

	result, err := deploymentsClient.Create(d)
	if err != nil {
		return errors.Wrap(err, "Couldn't create deployment client")
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// Remove application from Kubernetes on shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			deletePolicy := metav1.DeletePropagationForeground
			if err := deploymentsClient.Delete("demo-deployment", &metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				panic(err)
			}

			os.Exit(1)
		}
	}()

	return nil
}

func main() {
	var kubeconfig *string

	deploy := flag.Bool("kubernetes", false, "Deploy to kubernetes")
	replicas := flag.Int("replicas", 1, "Amount of instances to deploy")

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	if *deploy {
		if err := deployKubernetes(*kubeconfig, int32(*replicas)); err != nil {
			fmt.Printf("%v", err)
		}
	}

	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
