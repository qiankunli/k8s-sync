package kube

import (
	"k8s-sync/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// GetKubernetesClient returns the client if its possible in cluster, otherwise tries to read HOME
func GetKubernetesClient(config *config.Config) (*kubernetes.Clientset, error) {
	restConfig, err := GetKubernetesConfig(config)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

func GetKubernetesConfig(config *config.Config) (*rest.Config, error) {

	restClientConfig := &rest.Config{}
	restClientConfig.Host = config.Host
	restClientConfig.BearerToken = config.BearerToken
	restClientConfig.Burst = 1e6
	restClientConfig.QPS = 1e6
	restClientConfig.ContentType = "application/vnd.kubernetes.protobuf"

	return restClientConfig, nil
}
