// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"os"
	"regexp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	mangledRegExp = regexp.MustCompile("[^a-zA-Z0-9-_.]")
)

// GetKubeClientset returns a client object for a provided kubeconfig filepath
// if one is provided, or which uses the service account kubernetes gives to
// pods otherwise
func GetKubeClientset(kubeconfigPath string) (*kubernetes.Clientset, error) {
	fmt.Fprintf(os.Stdout, "Creating kube client object for kubeconfigPath %s\n", kubeconfigPath)

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		if kubeconfigPath != "" {
			return nil, fmt.Errorf("unable to load kubeconfig from %s: %w", kubeconfigPath, err)
		}

		return nil, fmt.Errorf("unable to load in-cluster config: %w", err)
	}

	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %w", err)
	}

	return client, nil
}

func MangleName(name string) string {
	// TODO need sensible unique naming scheme for deployments and secrets!
	return mangledRegExp.ReplaceAllString(name, "-")[:63]
}
