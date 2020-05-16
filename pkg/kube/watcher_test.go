package kube

import "testing"

func TestPodWatcher_Start(t *testing.T) {
	kubeClient, err := GetKubernetesClient("test")
	if err != nil {
		t.Error(err)
		return
	}
	w := NewPodWatcher(kubeClient)
	w.Start()
	stopper := make(chan int)
	<- stopper
}
