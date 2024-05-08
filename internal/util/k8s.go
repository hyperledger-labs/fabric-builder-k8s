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
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/watch"
	applycorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	typedBatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	watchtools "k8s.io/client-go/tools/watch"
	"k8s.io/utils/ptr"
)

const (
	namespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	jobTTL        = 5 * time.Minute

	ObjectNameSuffixLength int = 5

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

func waitForJob(
	ctx context.Context,
	client cache.Getter,
	jobName, namespace string,
	conditionFunc watchtools.ConditionFunc,
) (*batchv1.JobStatus, error) {
	fieldSelector := fields.OneTermEqualSelector("metadata.name", jobName)
	listWatch := cache.NewListWatchFromClient(client, "jobs", namespace, fieldSelector)

	ctx, cancel := watchtools.ContextWithOptionalTimeout(ctx, 0)
	defer cancel()

	event, err := watchtools.UntilWithSync(ctx, listWatch, &batchv1.Job{}, nil, conditionFunc)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, fmt.Errorf("no events received for job %s/%s", namespace, jobName)
	}

	job, ok := event.Object.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("event contained unexpected object %T while watching job %s/%s", job, namespace, jobName)
	}

	return &job.Status, nil
}

func waitForJobTermination(
	ctx context.Context,
	client cache.Getter,
	jobName, namespace string,
) (*batchv1.JobStatus, error) {
	jobTerminationCondition := func(event watch.Event) (bool, error) {
		if event.Type == watch.Deleted {
			return true, nil
		}

		job, ok := event.Object.(*batchv1.Job)
		if !ok {
			return false, fmt.Errorf(
				"event contained unexpected object %T while watching job %s/%s",
				job,
				namespace,
				jobName,
			)
		}

		for _, c := range job.Status.Conditions {
			if c.Type == batchv1.JobComplete && c.Status == "True" {
				return true, nil
			} else if c.Type == batchv1.JobFailed && c.Status == "True" {
				return true, fmt.Errorf("job %s/%s failed for reason %s: %s", namespace, jobName, c.Reason, c.Message)
			}
		}

		return false, nil
	}

	return waitForJob(ctx, client, jobName, namespace, jobTerminationCondition)
}

func WaitForChaincodeJob(
	ctx context.Context,
	logger *log.CmdLogger,
	client cache.Getter,
	job *batchv1.Job,
	chaincodeID string,
) error {
	logger.Debugf("Waiting for job %s/%s to terminate for chaincode ID %s", job.Namespace, job.Name, chaincodeID)

	_, err := waitForJobTermination(ctx, client, job.Name, job.Namespace)
	if err != nil {
		return fmt.Errorf(
			"error waiting for chaincode job %s/%s to terminate for chaincode ID %s: %w",
			job.Namespace,
			job.Name,
			chaincodeID,
			err,
		)
	}

	return fmt.Errorf(
		"chaincode job %s/%s for chaincode ID %s terminated",
		job.Namespace,
		job.Name,
		chaincodeID,
	)
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

func getChaincodeJobSpec(
	imageData *ImageJSON,
	namespace, serviceAccount, objectName, peerID string,
	chaincodeData *ChaincodeJSON,
) (*batchv1.Job, error) {
	chaincodeImage := imageData.Name + "@" + imageData.Digest

	jobName := objectName + "-" + rand.String(ObjectNameSuffixLength)

	labels, err := getLabels(chaincodeData)
	if err != nil {
		return nil, fmt.Errorf("error getting chaincode job labels for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	annotations := getAnnotations(peerID, chaincodeData)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        jobName,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: apiv1.PodSpec{
					ServiceAccountName: serviceAccount,
					Containers: []apiv1.Container{
						{
							Name:  "chaincode",
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
									SecretName: objectName,
								},
							},
						},
					},
				},
			},
			BackoffLimit:            ptr.To[int32](0),
			TTLSecondsAfterFinished: ptr.To[int32](int32(jobTTL / time.Second)),
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
		return fmt.Errorf("error applying chaincode secret definition for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	logger.Debugf(
		"Applied secrets for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		result.Namespace,
		result.Name,
	)

	return nil
}

func CreateChaincodeJob(
	ctx context.Context,
	logger *log.CmdLogger,
	jobsClient typedBatchv1.JobInterface,
	objectName, namespace, serviceAccount, peerID string,
	chaincodeData *ChaincodeJSON,
	imageData *ImageJSON,
) (*batchv1.Job, error) {
	jobDefinition, err := getChaincodeJobSpec(
		imageData,
		namespace,
		serviceAccount,
		objectName,
		peerID,
		chaincodeData,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting chaincode job definition for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	jobName := jobDefinition.ObjectMeta.Name

	logger.Debugf(
		"Creating chaincode job for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		namespace,
		jobName,
	)

	job, err := jobsClient.Create(ctx, jobDefinition, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf(
			"error creating chaincode job %s/%s for chaincode ID %s: %w",
			namespace,
			objectName,
			chaincodeData.ChaincodeID,
			err,
		)
	}

	logger.Debugf(
		"Created chaincode job for chaincode ID %s: %s/%s",
		chaincodeData.ChaincodeID,
		job.Namespace,
		job.Name,
	)

	return job, nil
}

// GetValidRfc1035LabelName returns a valid RFC 1035 label name with the format
// <prefix>-<truncated_chaincode_label>-<chaincode_run_hash> and space for a suffix if required.
func GetValidRfc1035LabelName(prefix, peerID string, chaincodeData *ChaincodeJSON, suffixLen int) string {
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
	runHashString := strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(runHash.Sum(nil)))

	// Remove unsafe characters from the chaincode package label
	packageID := NewChaincodePackageID(chaincodeData.ChaincodeID)
	re := regexp.MustCompile("[^-0-9a-z]")
	safeLabel := re.ReplaceAllString(strings.ToLower(packageID.Label), "")

	// Make sure the chaincode package label fits in the space available,
	// taking in to account the prefix, runHashString, two '-' separators,
	// and any required space for a suffix
	maxLabelLength := maxRfc1035LabelLength - len(prefix) - len(runHashString) - labelSeparators - suffixLen
	if maxLabelLength < len(safeLabel) {
		safeLabel = safeLabel[:maxLabelLength]
	}

	return prefix + "-" + safeLabel + "-" + runHashString
}
