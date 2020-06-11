package controller

import (
	"context"
	"fmt"
	"github.com/etiennecoutaud/curlme-controller/internal/curl"
	"github.com/etiennecoutaud/curlme-controller/internal/metrics"
	"github.com/etiennecoutaud/curlme-controller/internal/utils"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	cmInformer "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	cmListers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	samplescheme "k8s.io/sample-controller/pkg/generated/clientset/versioned/scheme"
)

const controllerAgentName = "configmap-curlme-controller"

// Event messages
const (
	SuccessSynced             = "Synced"
	MessageResourceSynced     = "configMap synced successfully"
	ErrCallingURL             = "Fail to retrieve message from URL"
	ErrParsingAnnotationValue = "Cannot parse x-k8s.io/curl-me-that annotation value"
)

// Controller struct for curlme controller
type Controller struct {
	kubeclientset   kubernetes.Interface
	configMapLister cmListers.ConfigMapLister
	configMapSynced cache.InformerSynced
	workqueue       workqueue.RateLimitingInterface
	recorder        record.EventRecorder
	curl            *curl.Curl
	metrics         *metrics.CurlmeMetrics
}

// New initiate a new controller
func New(
	kubeclientset kubernetes.Interface,
	configMapInformer cmInformer.ConfigMapInformer) *Controller {
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	c := curl.New()
	metrics := metrics.New()

	controller := &Controller{
		kubeclientset:   kubeclientset,
		configMapLister: configMapInformer.Lister(),
		configMapSynced: configMapInformer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ConfigMaps"),
		recorder:        recorder,
		curl:            c,
		metrics:         metrics,
	}

	klog.Info("Setting up event handlers")
	configMapInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newCm := new.(*corev1.ConfigMap)
			oldCm := old.(*corev1.ConfigMap)
			if utils.CompareAnnotation(oldCm.GetAnnotations(), newCm.GetAnnotations()) {
				// If CM is updated but annotation is the same, dont resync
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run start controller reconciliation loop
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting ConfigMap CurlMe controller")
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.configMapSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		c.workqueue.Forget(obj)
		c.metrics.CmSyncedCount.Inc()
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	cm, err := c.configMapLister.ConfigMaps(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("cm '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}
	klog.Infof("Config Map in handler : %s in %s namespace", cm.GetName(), cm.GetNamespace())

	key, url, err := utils.SplitAnnotationValue(utils.GetAnnotationValue(cm.GetAnnotations()))
	if err != nil {
		c.recorder.Event(cm, corev1.EventTypeWarning, ErrParsingAnnotationValue, err.Error())
		return nil
	}
	value, err := c.curl.CallingURL(url)
	if err != nil {
		c.recorder.Event(cm, corev1.EventTypeWarning, ErrCallingURL, err.Error())
		return nil
	}

	cmCopy := cm.DeepCopy()
	if len(cmCopy.Data) == 0 {
		cmCopy.Data = map[string]string{
			key: value,
		}
	} else {
		cmCopy.Data[key] = value
	}

	_, err = c.kubeclientset.CoreV1().ConfigMaps(namespace).Update(context.TODO(), cmCopy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	c.recorder.Event(cm, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) enqueueCm(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())

	if utils.ContainsAnnotation(object) {
		cm, err := c.configMapLister.ConfigMaps(object.GetNamespace()).Get(object.GetName())
		if err != nil {
			klog.V(4).Infof("cannot get configMap object %s in %s namespace", object.GetName(), object.GetNamespace())
			return
		}
		c.enqueueCm(cm)
		return
	}
}
