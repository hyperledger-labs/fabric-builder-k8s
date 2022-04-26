// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClientset returns a client object for a provided kubeconfig filepath
// if one is provided, or which uses the service account kubernetes gives to
// pods otherwise
func GetKubeClientset(kubeconfigPath string) (*kubernetes.Clientset, error) {

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		if kubeconfigPath != "" {
			return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)
		}

		return nil, fmt.Errorf("unable to load in-cluster config: %v", err)
	}

	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}

	return client, nil
}
