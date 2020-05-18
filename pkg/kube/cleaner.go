package kube

import (
	"github.com/rs/zerolog/log"
	"k8s-sync/pkg/config"
	"k8s-sync/pkg/db"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

type Cleaner struct {
	client          *kubernetes.Clientset
	batch           int32
	intervalSeconds time.Duration
	stopCh          chan bool
}

func NewCleaner(client *kubernetes.Clientset, config *config.Config) *Cleaner {
	cleaner := &Cleaner{
		client:          client,
		batch:           config.Batch,
		intervalSeconds: (time.Duration)(config.IntervalSeconds),
		stopCh:          make(chan bool),
	}
	return cleaner
}

func (c *Cleaner) Run() {
	// 可能会给apiserver 带来较大的负载
	for range time.Tick(c.intervalSeconds * time.Second) {
		select {
		case <-c.stopCh:
			return
		default:
		}
		db.Traverse(c.batch, func(pod db.PersistentPod) error {
			return c.TryClean(pod.Name)
		})
	}
}

func (c *Cleaner) TryClean(podName string) error {
	_, err := c.client.CoreV1().Pods("default").Get(podName, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Debug().Msgf("can not find pod name %s in k8s,delete it", podName)
			db.DeleteByName(podName)
			return nil
		}
		return err
	}
	return nil
}

func (c *Cleaner) Start() {
	go c.Run()
}

func (c *Cleaner) Stop() {
	c.stopCh <- true
	close(c.stopCh)
}
