// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/log"
	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Run struct {
	BuildOutputDirectory string
	RunMetadataDirectory string
	PeerID               string
	KubeconfigPath       string
	KubeNamespace        string
}

func (r *Run) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Running chaincode...")

	imageJsonPath := filepath.Join(r.BuildOutputDirectory, "/image.json")
	logger.Debugf("Reading %s...", imageJsonPath)

	_, err := os.Stat(imageJsonPath)
	if err != nil {
		return fmt.Errorf("unable to access %s: %w", imageJsonPath, err)
	}

	imageJsonContents, err := ioutil.ReadFile(imageJsonPath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", imageJsonPath, err)
	}

	var imageData util.ImageJson
	if err := json.Unmarshal(imageJsonContents, &imageData); err != nil {
		return fmt.Errorf("unable to process %s: %w", imageJsonPath, err)
	}

	logger.Debugf("Image name: %s\nImage digest: %s\n", imageData.Name, imageData.Digest)

	chaincodeJsonPath := filepath.Join(r.RunMetadataDirectory, "/chaincode.json")
	logger.Debugf("Reading %s...", chaincodeJsonPath)

	_, err = os.Stat(chaincodeJsonPath)
	if err != nil {
		return fmt.Errorf("unable to access %s: %w", chaincodeJsonPath, err)
	}

	chaincodeJsonContents, err := ioutil.ReadFile(chaincodeJsonPath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", chaincodeJsonPath, err)
	}

	var chaincodeData util.ChaincodeJson
	if err := json.Unmarshal(chaincodeJsonContents, &chaincodeData); err != nil {
		return fmt.Errorf("unable to process %s: %w", chaincodeJsonPath, err)
	}

	logger.Debugf("Chaincode ID: %s\n", chaincodeData.ChaincodeID)

	clientset, err := util.GetKubeClientset(logger, r.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("unable to connect kubernetes client for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}

	secretsClient := clientset.CoreV1().Secrets(r.KubeNamespace)

	secret := util.GetChaincodeSecretApplyConfiguration(r.KubeNamespace, r.PeerID, chaincodeData)

	s, err := secretsClient.Apply(ctx, secret, metav1.ApplyOptions{FieldManager: "fabric-builder-k8s"})
	if err != nil {
		return fmt.Errorf("unable to create kubernetes secret for chaincode ID %s: %w", chaincodeData.ChaincodeID, err)
	}
	logger.Debugf("Applied secret %s\n", s.Name)

	podsClient := clientset.CoreV1().Pods(r.KubeNamespace)

	podName := util.GetPodName(chaincodeData.MspID, r.PeerID, chaincodeData.ChaincodeID)

	pod := util.GetChaincodePodObject(imageData, r.KubeNamespace, podName, r.PeerID, chaincodeData)

	createAttempts := 0
	for {
		createAttempts += 1
		p, err := podsClient.Create(ctx, pod, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				if createAttempts > 3 {
					// give up
					return fmt.Errorf("unable to create chaincode pod %s/%s for chaincode ID %s on final attempt: %w", r.KubeNamespace, podName, chaincodeData.ChaincodeID, err)
				}

				err = podsClient.Delete(ctx, podName, metav1.DeleteOptions{})
				if err != nil {
					if !errors.IsNotFound(err) {
						logger.Printf("Error deleting existing chaincode pod for chaincode ID %s: %v", chaincodeData.ChaincodeID, err)
					}
				}

				_, err := util.WaitForPodTermination(ctx, time.Minute, podsClient, podName, r.KubeNamespace)
				if err != nil {
					if !errors.IsNotFound(err) {
						logger.Printf("Error waiting for existing chaincode pod to terminate for chaincode ID %s: %v", chaincodeData.ChaincodeID, err)
					}
				}

				// try again
				continue
			}

			return fmt.Errorf("unable to create chaincode pod %s/%s for chaincode ID %s: %w", r.KubeNamespace, podName, chaincodeData.ChaincodeID, err)
		}

		logger.Debugf("Created chaincode pod for chaincode ID %s: %s/%s", chaincodeData.ChaincodeID, p.Namespace, p.Name)
		break
	}

	_, err = util.WaitForPodRunning(ctx, time.Minute, podsClient, podName, r.KubeNamespace)
	if err != nil {
		return fmt.Errorf("error waiting for chaincode pod %s/%s for chaincode ID %s: %w", r.KubeNamespace, podName, chaincodeData.ChaincodeID, err)
	}

	status, err := util.WaitForPodTermination(ctx, 0, podsClient, podName, r.KubeNamespace)
	if err != nil {
		return fmt.Errorf("error waiting for chaincode pod %s/%s to terminate for chaincode ID %s: %w", r.KubeNamespace, podName, chaincodeData.ChaincodeID, err)
	}
	if status != nil {
		return fmt.Errorf("chaincode pod %s/%s for chaincode ID %s terminated %s: %s", r.KubeNamespace, podName, chaincodeData.ChaincodeID, status.Reason, status.Message)
	}

	return fmt.Errorf("unexpected chaincode pod termination for chaincode ID %s", chaincodeData.ChaincodeID)
}
