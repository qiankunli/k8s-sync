package kube

import (
	"k8s-sync/pkg/config"
	"k8s-sync/pkg/db"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/rs/zerolog/log"
)

type PodWatcher struct {
	informer cache.SharedInformer
	env      string
	stopper  chan struct {
}
}

func NewPodWatcher(client *kubernetes.Clientset, config *config.Config) *PodWatcher {
	factory := informers.NewSharedInformerFactory(client, 0)
	informer := factory.Core().V1().Pods().Informer()

	watcher := &PodWatcher{
		informer: informer,
		env:      config.Env,
		stopper:  make(chan struct{}),
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: watcher.OnAdd,
		// This invokes the ReplicaSet for every pod change, eg: host assignment. Though this might seem like
		// overkill the most frequent pod update is status, and the associated ReplicaSet will only list from
		// local storage, so it should be ok.
		UpdateFunc: watcher.OnUpdate,
		DeleteFunc: watcher.OnDelete,
	})

	return watcher
}

func (w *PodWatcher) OnAdd(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Debug().Msgf("onAdd pod: %s", pod.ObjectMeta.Name)
	if pod.DeletionTimestamp != nil {
		// on a restart of the controller manager, it's possible a new pod shows up in a state that
		// is already pending deletion. Prevent the pod from being a creation observation.
		w.OnDelete(pod)
		return
	}
	if err := db.SaveOrUpdate(db.NewPersistentPod(pod,w.env)); err != nil {
		log.Err(err)
	}
}

func (w *PodWatcher) OnUpdate(oldObj, newObj interface{}) {
	pod := newObj.(*v1.Pod)
	log.Debug().Msgf("onUpdate pod: %s", pod.ObjectMeta.Name)
	if err := db.UpdateByName(db.NewPersistentPod(pod,w.env)); err != nil {
		log.Err(err)
	}
}

func (w *PodWatcher) OnDelete(obj interface{}) {
	pod := obj.(*v1.Pod)
	log.Debug().Msgf("OnDelete pod: %s", pod.ObjectMeta.Name)
	if err := db.DeleteByName(pod.Name); err != nil {
		log.Err(err)
	}
}

func (w *PodWatcher) Start() {
	go w.informer.Run(w.stopper)
}

func (w *PodWatcher) Stop() {
	w.stopper <- struct{}{}
	close(w.stopper)
}
