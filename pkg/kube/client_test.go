package kube

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetKubernetesClient(t *testing.T) {

	kubeClient, err := GetKubernetesClient("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	pod, err := kubeClient.CoreV1().Pods("default").Get("fm-abtest-backend-stable-5c45886f8c-z45bq", metav1.GetOptions{})
	fmt.Println(err)
	fmt.Println(pod)
	fmt.Println(pod.ObjectMeta.Name)
	fmt.Println(pod.Labels["isolation"])
	fmt.Println(pod.Labels["appName"])
	fmt.Println(pod.Status.HostIP)
	fmt.Println(pod.Status.Phase)


}
