package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func terminateInstance(client *kubernetes.Clientset, pod v1.Pod) {
	gracePeriod := int64(10)
	deleteopts := &meta_v1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	}
	client.Core().Pods(apiv1.NamespaceDefault).Delete(pod.Name, deleteopts)
}

func scheduleNextTermination(client *kubernetes.Clientset) (v1.Pod, error) {
	list, err := client.CoreV1().Pods(apiv1.NamespaceDefault).List(apiv1.ListOptions{})
	if err != nil {
		return v1.Pod{}, errors.Wrap(err, "Couldn't get list of pods.")
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	i := r.Intn(len(list.Items))

	time.Sleep(10 * time.Second)
	return list.Items[i], nil
}

func configureClient() (*kubernetes.Clientset, error) {
	home := homeDir()
	configPath := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	config, err := clientcmd.BuildConfigFromFlags("", *configPath)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't build configuration")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't build client with given config")
	}

	return clientset, nil
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	terminate := make(chan v1.Pod, 10)
	defer close(terminate)

	client, err := configureClient()
	if err != nil {
		panic(err)
	}

	for {
		go func() {
			pod, err := scheduleNextTermination(client)
			if err != nil {
				panic(err)
			}
			terminate <- pod
		}()

		select {
		case <-c:
			fmt.Println("Shutting down")
			os.Exit(1)
		case pod := <-terminate:
			fmt.Printf("Terminating pod: %s \n", pod.Name)
			terminateInstance(client, pod)
		}
	}
}
