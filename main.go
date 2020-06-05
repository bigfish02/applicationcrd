package main

import (
	"flag"
	"os"
	"path"
	"time"

	"github.com/bigfish02/applicationcrd/pkg/signals"

	clientset "github.com/bigfish02/applicationcrd/pkg/generated/clientset/versioned"
	crdinformers "github.com/bigfish02/applicationcrd/pkg/generated/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func main() {
	home := os.Getenv("HOME")
	kubeconfig := flag.String("kubeconfig", path.Join(home, "./.kube", "config"), "Path to kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		klog.Fatalf("Error Building kubeconfig: %s\n", err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error New kubernetes client: %s\n", err.Error())
	}
	crdClient, err := clientset.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Error New application clientset: %s\n", err.Error())
	}
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	crdInformerFactory := crdinformers.NewSharedInformerFactory(crdClient, time.Second*30)
	controller := NewController(kubeClient, crdClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		kubeInformerFactory.Core().V1().Services(),
		crdInformerFactory.Xiaohongshu().V1().Applications())

	stopCh := signals.SetupSignalHandler()
	go kubeInformerFactory.Start(stopCh)
	go crdInformerFactory.Start(stopCh)
	if err = controller.Run(2, stopCh); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}
}
