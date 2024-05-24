// SPDX-License-Identifier: Apache-2.0

package util

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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

	// Defaults.
	DefaultNamespace          string = "default"
	DefaultObjectNamePrefix   string = "hlfcc"
	DefaultServiceAccountName string = "default"

	// Mutual TLS auth client key and cert paths in the chaincode container.
	TLSClientKeyPath      string = "/etc/hyperledger/fabric/client.key"
	TLSClientCertPath     string = "/etc/hyperledger/fabric/client.crt"
	TLSClientKeyFile      string = "/etc/hyperledger/fabric/client_pem.key"
	TLSClientCertFile     string = "/etc/hyperledger/fabric/client_pem.crt"
	TLSClientRootCertFile string = "/etc/hyperledger/fabric/peer.crt"
)

func waitForPod(
	ctx context.Context,
	timeout time.Duration,
	podsClient v1.PodInterface,
	podName, namespace string,
	conditionFunc watchtools.ConditionFunc,
) (*apiv1.PodStatus, error) {
	fieldSelector := fields.OneTermEqualSelector("metadata.name", podName).String()

	listWatch := &cache.ListWatch{
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
	// client to give it to avoid forbidden errors for pod list.
	// var client kubernetes.Interface
	// listWatch := cache.NewListWatchFromClient(client, "pods", namespace, fieldSelector)

	ctx, cancel := watchtools.ContextWithOptionalTimeout(ctx, timeout)
	defer cancel()

	event, err := watchtools.UntilWithSync(ctx, listWatch, &apiv1.Pod{}, nil, conditionFunc)
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

	return &pod.Status, nil
}

func waitForPodRunning(
	ctx context.Context,
	timeout time.Duration,
	podsClient v1.PodInterface,
	podName, namespace string,
) (*apiv1.PodStatus, error) {
	podRunningCondition := func(event watch.Event) (bool, error) {
		pod, ok := event.Object.(*apiv1.Pod)
		if !ok {
			return false, fmt.Errorf(
				"unexpected object while watching pod %s/%s",
				namespace,
				podName,
			)
		}

		phase := pod.Status.Phase
		if phase == apiv1.PodRunning {
			return true, nil
		}

		return false, nil
	}

	return waitForPod(ctx, timeout, podsClient, podName, namespace, podRunningCondition)
}

func waitForPodTermination(
	ctx context.Context,
	timeout time.Duration,
	podsClient v1.PodInterface,
	podName, namespace string,
) (*apiv1.PodStatus, error) {
	podTerminationCondition := func(event watch.Event) (bool, error) {
		if event.Type == watch.Deleted {
			return true, nil
		}

		pod, ok := event.Object.(*apiv1.Pod)
		if !ok {
			return false, fmt.Errorf(
				"unexpected object while watching pod %s/%s",
				namespace,
				podName,
			)
		}

		phase := pod.Status.Phase
		if phase != apiv1.PodRunning {
			return true, nil
		}

		return false, nil
	}

	return waitForPod(ctx, timeout, podsClient, podName, namespace, podTerminationCondition)
}

func WaitForChaincodePod(
	ctx context.Context,
	logger *log.CmdLogger,
	podsClient v1.PodInterface,
	pod *apiv1.Pod,
	chaincodeID string,
) error {
	logger.Debugf("Waiting for pod %s/%s for chaincode ID %s", pod.Namespace, pod.Name, chaincodeID)

	_, err := waitForPodRunning(ctx, time.Minute, podsClient, pod.Name, pod.Namespace)
	if err != nil {
		return fmt.Errorf(
			"error waiting for chaincode pod %s/%s for chaincode ID %s: %w",
			pod.Namespace,
			pod.Name,
			chaincodeID,
			err,
		)
	}

	status, err := waitForPodTermination(ctx, 0, podsClient, pod.Name, pod.Namespace)
	if err != nil {
		return fmt.Errorf(
			"error waiting for chaincode pod %s/%s to terminate for chaincode ID %s: %w",
			pod.Namespace,
			pod.Name,
			chaincodeID,
			err,
		)
	}

	if status != nil {
		return fmt.Errorf(
			"chaincode pod %s/%s for chaincode ID %s terminated %s: %s",
			pod.Namespace,
			pod.Name,
			chaincodeID,
			status.Reason,
			status.Message,
		)
	}

	return fmt.Errorf("unexpected chaincode pod termination for chaincode ID %s", chaincodeID)
}

// GetKubeClientset returns a client object for a provided kubeconfig filepath
// if one is provided, or which uses the service account kubernetes gives to
// pods otherwise.
func GetKubeClientset(logger *log.CmdLogger, kubeconfigPath string) (*kubernetes.Clientset, error) {
	logger.Debugf("Creating kube client object for kubeconfigPath %s", kubeconfigPath)

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
	namespace, err := os.ReadFile(namespacePath)
	if err != nil {
		return "", fmt.Errorf("unable to read namespace from %s: %w", namespacePath, err)
	}

	return string(namespace), nil
}

func getLabels(chaincodeData *ChaincodeJSON) (map[string]string, error) {
	packageID := NewChaincodePackageID(chaincodeData.ChaincodeID)

	packageHashBytes, err := hex.DecodeString(packageID.Hash)
	if err != nil {
		return nil, fmt.Errorf("error decoding chaincode package hash %s: %w", packageID.Hash, err)
	}

	encodedPackageHash := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(packageHashBytes)

	return map[string]string{
		"app.kubernetes.io/name":       "hyperledger-fabric",
		"app.kubernetes.io/component":  "chaincode",
		"app.kubernetes.io/created-by": "fabric-builder-k8s",
		"app.kubernetes.io/managed-by": "fabric-builder-k8s",
		"fabric-builder-k8s-cclabel":   packageID.Label,
		"fabric-builder-k8s-cchash":    encodedPackageHash,
	}, nil
}

func getAnnotations(peerID string, chaincodeData *ChaincodeJSON) map[string]string {
	return map[string]string{
		"fabric-builder-k8s-ccid":        chaincodeData.ChaincodeID,
		"fabric-builder-k8s-mspid":       chaincodeData.MspID,
		"fabric-builder-k8s-peeraddress": chaincodeData.PeerAddress,
		"fabric-builder-k8s-peerid":      peerID,
	}
}

func getChaincodePodObject(
	imageData *ImageJSON,
	namespace, serviceAccount, podName, peerID string,
	chaincodeData *ChaincodeJSON,
) (*apiv1.Pod, error) {
	chaincodeImage := imageData.Name + "@" + imageData.Digest

	labels, err := getLabels(chaincodeData)
	if err != nil {
		return nil, fmt.Errorf("error getting chaincode job labels for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	annotations := getAnnotations(peerID, chaincodeData)

	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        podName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: apiv1.PodSpec{
			ServiceAccountName: serviceAccount,
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
			RestartPolicy: apiv1.RestartPolicyNever,
			Volumes: []apiv1.Volume{
				{
					Name: "certs",
					VolumeSource: apiv1.VolumeSource{
						Secret: &apiv1.SecretVolumeSource{
							SecretName: podName,
						},
					},
				},
			},
		},
	}, nil
}

func getChaincodeSecretApplyConfiguration(
	secretName, namespace, peerID string,
	chaincodeData *ChaincodeJSON,
) (*applycorev1.SecretApplyConfiguration, error) {
	labels, err := getLabels(chaincodeData)
	if err != nil {
		return nil, fmt.Errorf("error getting chaincode job labels for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	annotations := getAnnotations(peerID, chaincodeData)

	data := map[string]string{
		"peer.crt":       chaincodeData.RootCert,
		"client_pem.crt": chaincodeData.ClientCert,
		"client_pem.key": chaincodeData.ClientKey,
		"client.crt":     base64.StdEncoding.EncodeToString([]byte(chaincodeData.ClientCert)),
		"client.key":     base64.StdEncoding.EncodeToString([]byte(chaincodeData.ClientKey)),
	}

	return applycorev1.
		Secret(secretName, namespace).
		WithAnnotations(annotations).
		WithLabels(labels).
		WithStringData(data).
		WithType(apiv1.SecretTypeOpaque), nil
}

func ApplyChaincodeSecrets(
	ctx context.Context,
	logger *log.CmdLogger,
	secretsClient v1.SecretInterface,
	secretName, namespace, peerID string,
	chaincodeData *ChaincodeJSON,
) error {
	secret, err := getChaincodeSecretApplyConfiguration(secretName, namespace, peerID, chaincodeData)
	if err != nil {
		return fmt.Errorf("error getting chaincode secret definition for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	result, err := secretsClient.Apply(
		ctx,
		secret,
		metav1.ApplyOptions{FieldManager: "fabric-builder-k8s"},
	)
	if err != nil {
		return err
	}

	logger.Debugf(
		"Applied secrets for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		result.Namespace,
		result.Name,
	)

	return nil
}

func deleteChaincodePod(
	ctx context.Context,
	logger *log.CmdLogger,
	podsClient v1.PodInterface,
	podName, namespace string,
	chaincodeData *ChaincodeJSON,
) error {
	logger.Debugf(
		"Deleting any existing chaincode pod for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		namespace,
		podName,
	)

	err := podsClient.Delete(ctx, podName, metav1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Debugf(
				"No existing chaincode pod for chaincode ID %s: %s/%s",
				chaincodeData.ChaincodeID,
				namespace,
				podName,
			)

			return nil
		}

		return err
	}

	logger.Debugf(
		"Waiting for existing chaincode pod to terminate for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		namespace,
		podName,
	)

	_, err = waitForPodTermination(ctx, time.Minute, podsClient, podName, namespace)
	if err != nil {
		return err
	}

	logger.Debugf(
		"Existing chaincode pod deleted for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		namespace,
		podName,
	)

	return nil
}

func CreateChaincodePod(
	ctx context.Context,
	logger *log.CmdLogger,
	podsClient v1.PodInterface,
	objectName, namespace, serviceAccount, peerID string,
	chaincodeData *ChaincodeJSON,
	imageData *ImageJSON,
) (*apiv1.Pod, error) {
	podDefinition, err := getChaincodePodObject(
		imageData,
		namespace,
		serviceAccount,
		objectName,
		peerID,
		chaincodeData,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting chaincode pod definition for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	err = deleteChaincodePod(ctx, logger, podsClient, objectName, namespace, chaincodeData)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to delete existing chaincode pod %s/%s for chaincode ID %s: %w",
			namespace,
			objectName,
			chaincodeData.ChaincodeID,
			err,
		)
	}

	logger.Debugf(
		"Creating chaincode pod for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		namespace,
		objectName,
	)

	pod, err := podsClient.Create(ctx, podDefinition, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf(
			"unable to create chaincode pod %s/%s for chaincode ID %s: %w",
			namespace,
			objectName,
			chaincodeData.ChaincodeID,
			err,
		)
	}

	logger.Debugf(
		"Created chaincode pod for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		pod.Namespace,
		pod.Name,
	)

	return pod, nil
}

// GetValidRfc1035LabelName returns a valid RFC 1035 label name with the format
// <prefix>-<truncated_chaincode_label>-<chaincode_run_hash>.
func GetValidRfc1035LabelName(prefix, peerID string, chaincodeData *ChaincodeJSON) string {
	const (
		maxRfc1035LabelLength = 63
		labelSeparators       = 2
	)

	runHash := fnv.New64a()
	runHash.Write([]byte(prefix))
	runHash.Write([]byte(peerID))
	runHash.Write([]byte(chaincodeData.PeerAddress))
	runHash.Write([]byte(chaincodeData.MspID))
	runHash.Write([]byte(chaincodeData.ChaincodeID))
	suffix := strings.ToLower(base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(runHash.Sum(nil)))

	// Remove unsafe characters from the chaincode package label
	packageID := NewChaincodePackageID(chaincodeData.ChaincodeID)
	re := regexp.MustCompile("[^-0-9a-z]")
	safeLabel := re.ReplaceAllString(strings.ToLower(packageID.Label), "")

	// Make sure the chaincode package label fits in the space available,
	// taking in to account the prefix, suffix, and two '-' separators
	maxLabelLength := maxRfc1035LabelLength - len(prefix) - len(suffix) - labelSeparators
	if maxLabelLength < len(safeLabel) {
		safeLabel = safeLabel[:maxLabelLength]
	}

	return prefix + "-" + safeLabel + "-" + suffix
}
