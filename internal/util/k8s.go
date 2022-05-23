// SPDX-License-Identifier: Apache-2.0

package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	watchtools "k8s.io/client-go/tools/watch"
)

const (
	namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

	// Mutual TLS auth client key and cert paths in the chaincode container
	TLSClientKeyPath      string = "/etc/hyperledger/fabric/client.key"
	TLSClientCertPath     string = "/etc/hyperledger/fabric/client.crt"
	TLSClientKeyFile      string = "/etc/hyperledger/fabric/client_pem.key"
	TLSClientCertFile     string = "/etc/hyperledger/fabric/client_pem.crt"
	TLSClientRootCertFile string = "/etc/hyperledger/fabric/peer.crt"
)

var (
	mangledRegExp = regexp.MustCompile("[^a-zA-Z0-9-_.]")
)

func waitForPod(ctx context.Context, timeout time.Duration, podsClient v1.PodInterface, podName, namespace string, conditionFunc watchtools.ConditionFunc) (*apiv1.PodStatus, error) {
	fieldSelector := fields.OneTermEqualSelector("metadata.name", podName).String()

	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fieldSelector
			return podsClient.List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fieldSelector
			return podsClient.Watch(context.TODO(), options)
		},
	}

	// TODO it might be nice to use NewListWatchFromClient instead but not sure what
	// client to give it to avoid forbidden errors for pod list
	// var client kubernetes.Interface
	// lw := cache.NewListWatchFromClient(client, "pods", namespace, fieldSelector)

	ctx, cancel := watchtools.ContextWithOptionalTimeout(ctx, timeout)
	defer cancel()

	event, err := watchtools.UntilWithSync(ctx, lw, &apiv1.Pod{}, nil, conditionFunc)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("no events received for pod %s/%s", namespace, podName)
	}

	pod, ok := event.Object.(*apiv1.Pod)
	if !ok {
		return nil, fmt.Errorf("unexpected object while watching pod %s/%s", namespace, podName)
	}
	status := pod.Status

	return &status, nil
}

func WaitForPodRunning(ctx context.Context, timeout time.Duration, podsClient v1.PodInterface, podName, namespace string) (*apiv1.PodStatus, error) {
	podRunningCondition := func(event watch.Event) (bool, error) {
		pod, ok := event.Object.(*apiv1.Pod)
		if !ok {
			return false, fmt.Errorf("unexpected object while watching pod %s/%s", namespace, podName)
		}

		phase := pod.Status.Phase
		if phase == apiv1.PodRunning {
			return true, nil
		}

		return false, nil
	}

	return waitForPod(ctx, timeout, podsClient, podName, namespace, podRunningCondition)
}

func WaitForPodTermination(ctx context.Context, timeout time.Duration, podsClient v1.PodInterface, podName, namespace string) (*apiv1.PodStatus, error) {
	podTerminationCondition := func(event watch.Event) (bool, error) {
		switch event.Type {
		case watch.Deleted:
			return true, nil
		}

		pod, ok := event.Object.(*apiv1.Pod)
		if !ok {
			return false, fmt.Errorf("unexpected object while watching pod %s/%s", namespace, podName)
		}

		phase := pod.Status.Phase
		if phase != apiv1.PodRunning {
			return true, nil
		}

		return false, nil
	}

	return waitForPod(ctx, timeout, podsClient, podName, namespace, podTerminationCondition)
}

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

func GetKubeNamespace() (string, error) {
	namespace, err := ioutil.ReadFile(namespacePath)
	if err != nil {
		return "", fmt.Errorf("unable to read namespace from %s: %w", namespacePath, err)
	}

	return string(namespace), nil
}

func GetChaincodePodObject(imageData ImageJson, namespace, podName, peerID string, chaincodeData ChaincodeJson) *apiv1.Pod {
	chaincodeImage := imageData.Name + "@" + imageData.Digest

	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "fabric",
				"app.kubernetes.io/component":  "chaincode",
				"app.kubernetes.io/created-by": "fabric-builder-k8s",
				"app.kubernetes.io/managed-by": "fabric-builder-k8s",
				"fabric-builder-k8s-mspid":     chaincodeData.MspID,
				"fabric-builder-k8s-peerid":    peerID,
			},
			Annotations: map[string]string{
				"fabric-builder-k8s-ccid": chaincodeData.ChaincodeID,
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "main",
					Image: chaincodeImage,
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "certs",
							MountPath: "/etc/hyperledger/fabric",
							ReadOnly:  true,
						},
					},
					Env: []apiv1.EnvVar{
						{
							Name:  "CORE_CHAINCODE_ID_NAME",
							Value: chaincodeData.ChaincodeID,
						},
						{
							Name:  "CORE_PEER_ADDRESS",
							Value: chaincodeData.PeerAddress,
						},
						{
							Name:  "CORE_PEER_TLS_ENABLED",
							Value: "true", // TODO only if there are certs?
						},
						{
							Name:  "CORE_PEER_TLS_ROOTCERT_FILE",
							Value: TLSClientRootCertFile,
						},
						{
							Name:  "CORE_TLS_CLIENT_KEY_PATH",
							Value: TLSClientKeyPath,
						},
						{
							Name:  "CORE_TLS_CLIENT_CERT_PATH",
							Value: TLSClientCertPath,
						},
						{
							Name:  "CORE_TLS_CLIENT_KEY_FILE",
							Value: TLSClientKeyFile,
						},
						{
							Name:  "CORE_TLS_CLIENT_CERT_FILE",
							Value: TLSClientCertFile,
						},
						{
							Name:  "CORE_PEER_LOCALMSPID",
							Value: chaincodeData.MspID,
						},
					},
				},
			},
			Volumes: []apiv1.Volume{
				{
					Name: "certs",
					VolumeSource: apiv1.VolumeSource{
						Secret: &apiv1.SecretVolumeSource{
							SecretName: getSecretName(chaincodeData.MspID, peerID, chaincodeData.ChaincodeID),
						},
					},
				},
			},
		},
	}
}

func GetChaincodeSecretApplyConfiguration(namespace, peerID string, chaincodeData ChaincodeJson) *applycorev1.SecretApplyConfiguration {
	name := getSecretName(chaincodeData.MspID, peerID, chaincodeData.ChaincodeID)

	annotations := map[string]string{
		"fabric-builder-k8s-ccid": chaincodeData.ChaincodeID,
	}

	data := map[string]string{
		"peer.crt":       chaincodeData.RootCert,
		"client_pem.crt": chaincodeData.ClientCert,
		"client_pem.key": chaincodeData.ClientKey,
		"client.crt":     base64.StdEncoding.EncodeToString([]byte(chaincodeData.ClientCert)),
		"client.key":     base64.StdEncoding.EncodeToString([]byte(chaincodeData.ClientKey)),
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       "fabric",
		"app.kubernetes.io/component":  "chaincode",
		"app.kubernetes.io/created-by": "fabric-builder-k8s",
		"app.kubernetes.io/managed-by": "fabric-builder-k8s",
		"fabric-builder-k8s-mspid":     chaincodeData.MspID,
		"fabric-builder-k8s-peerid":    peerID,
	}

	return applycorev1.
		Secret(name, namespace).
		WithAnnotations(annotations).
		WithLabels(labels).
		WithStringData(data).
		WithType(apiv1.SecretTypeOpaque)
}

func GetChaincodeSecretObject(namespace, peerID string, chaincodeData ChaincodeJson) *apiv1.Secret {
	return &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getSecretName(chaincodeData.MspID, peerID, chaincodeData.ChaincodeID),
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       "fabric",
				"app.kubernetes.io/component":  "chaincode",
				"app.kubernetes.io/created-by": "fabric-builder-k8s",
				"app.kubernetes.io/managed-by": "fabric-builder-k8s",
				"fabric-builder-k8s-mspid":     chaincodeData.MspID,
				"fabric-builder-k8s-peerid":    peerID,
			},
			Annotations: map[string]string{
				"fabric-builder-k8s-ccid": chaincodeData.ChaincodeID,
			},
		},
		Type: apiv1.SecretTypeOpaque,
		StringData: map[string]string{
			"peer.crt":       chaincodeData.RootCert,
			"client_pem.crt": chaincodeData.ClientCert,
			"client_pem.key": chaincodeData.ClientKey,
			"client.crt":     base64.StdEncoding.EncodeToString([]byte(chaincodeData.ClientCert)),
			"client.key":     base64.StdEncoding.EncodeToString([]byte(chaincodeData.ClientKey)),
		},
	}
}

func GetPodName(mspID, peerID, chaincodeID string) string {
	return mangleName("cc-" + mspID + "-" + peerID + chaincodeID)
}

func getSecretName(mspID, peerID, chaincodeID string) string {
	return mangleName("cc-" + mspID + "-" + peerID + chaincodeID)
}

func mangleName(name string) string {
	// TODO need sensible unique naming scheme for deployments and secrets!
	return strings.ToLower(mangledRegExp.ReplaceAllString(name, "-")[:63])
}
