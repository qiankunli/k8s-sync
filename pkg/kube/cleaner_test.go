package kube

import (
	"fmt"
	"testing"
	"time"
)

func TestTryClean(t *testing.T) {
	kubeClient, err := GetKubernetesClient("test")
	if err != nil {
		t.Error(err)
		return
	}
	c := NewCleaner(kubeClient)
	err = c.TryClean("abc")
	fmt.Println(err)
}

func TestCleaner_Run(t *testing.T) {
	kubeClient, err := GetKubernetesClient("test")
	if err != nil {
		t.Error(err)
		return
	}
	c := NewCleaner(kubeClient)
	c.Start()
	time.Sleep(100000 * time.Second)
	fmt.Println("main stop1")
	c.Stop()
	fmt.Println("main stop2")
}
