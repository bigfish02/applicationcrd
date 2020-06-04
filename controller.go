package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	crdv1 "github.com/bigfish02/applicationcrd/pkg/apis/application/v1"
	clientset "github.com/bigfish02/applicationcrd/pkg/generated/clientset/versioned"
	"github.com/bigfish02/applicationcrd/pkg/generated/clientset/versioned/scheme"
	crdinformers "github.com/bigfish02/applicationcrd/pkg/generated/informers/externalversions/application/v1"
	crdlisters "github.com/bigfish02/applicationcrd/pkg/generated/listers/application/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

const controllerAgentName = "application-controller"

type Controller struct {
	kubeclientset        kubernetes.Interface
	applicationclientset clientset.Interface

	deploymentLister appslisters.DeploymentLister
	deploymentSynced cache.InformerSynced

	serviceLister corelisters.ServiceLister
	serviceSynced cache.InformerSynced

	applicationLister crdlisters.ApplicationLister
	applicationSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	recorder  record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	applicationclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	serviceInformer coreinformers.ServiceInformer,
	crdInformer crdinformers.ApplicationInformer) *Controller {

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{
		Interface: kubeclientset.CoreV1().Events(""),
	})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:        kubeclientset,
		applicationclientset: applicationclientset,
		deploymentLister:     deploymentInformer.Lister(),
		deploymentSynced:     deploymentInformer.Informer().HasSynced,
		serviceLister:        serviceInformer.Lister(),
		serviceSynced:        serviceInformer.Informer().HasSynced,
		applicationLister:    crdInformer.Lister(),
		applicationSynced:    crdInformer.Informer().HasSynced,
		workqueue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Applications"),
		recorder:             recorder,
	}

	crdInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.enqueueApplication,
		DeleteFunc: controller.enqueueApplicationForDelete,
		UpdateFunc: func(old, new interface{}) {
			oldApplication := old.(*crdv1.Application)
			newApplication := new.(*crdv1.Application)
			if oldApplication.ResourceVersion == newApplication.ResourceVersion {
				return
			}
			controller.enqueueApplication(new)
		},
	})

	return controller
}

func (c *Controller) enqueueApplication(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}
func (c *Controller) enqueueApplicationForDelete(obj interface{}) {
	var key string
	var err error
	if key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()
	klog.Info("Starting application controller loop")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentSynced, c.applicationSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	klog.Info("Start workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	<-stopCh
	klog.Info("Shuting down workers")
	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {

	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return true
	}
	if shutdown {
		return false
	}
	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	app, err := c.applicationLister.Applications(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("start deleting application: %s in namespace", name, namespace)
			return nil
		}
		return err
	}
	klog.Infof("start to process %v", app)
	return nil
}
