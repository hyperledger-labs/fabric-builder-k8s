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

func (r *Run) Run() error {
	imageJsonPath := filepath.Join(r.BuildOutputDirectory, "/image.json")
	chaincodeJsonPath := filepath.Join(r.RunMetadataDirectory, "/chaincode.json")

	fmt.Println("Reading image.json...")
	_, err := os.Stat(imageJsonPath)
	if err != nil {
		return fmt.Errorf("unable to access image.json: %w", err)
	}

	imageJsonContents, err := ioutil.ReadFile(imageJsonPath)
	if err != nil {
		return fmt.Errorf("unable to read image.json: %w", err)
	}

	var imageData util.ImageJson
	if err := json.Unmarshal(imageJsonContents, &imageData); err != nil {
		return fmt.Errorf("unable to process image.json: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Image name: %s\nImage digest: %s\n", imageData.Name, imageData.Digest)

	fmt.Println("Reading chaincode.json...")
	_, err = os.Stat(chaincodeJsonPath)
	if err != nil {
		return fmt.Errorf("unable to access chaincode.json: %w", err)
	}

	chaincodeJsonContents, err := ioutil.ReadFile(chaincodeJsonPath)
	if err != nil {
		return fmt.Errorf("unable to read chaincode.json: %w", err)
	}

	var chaincodeData util.ChaincodeJson
	if err := json.Unmarshal(chaincodeJsonContents, &chaincodeData); err != nil {
		return fmt.Errorf("unable to process chaincode.json: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Chaincode ID: %s\n", chaincodeData.ChaincodeID)

	clientset, err := util.GetKubeClientset(r.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("unable to connect kubernetes client: %w", err)
	}

	secretsClient := clientset.CoreV1().Secrets(r.KubeNamespace)

	secret := util.GetChaincodeSecretApplyConfiguration(r.KubeNamespace, r.PeerID, chaincodeData)

	s, err := secretsClient.Apply(context.TODO(), secret, metav1.ApplyOptions{FieldManager: "fabric-builder-k8s"})
	if err != nil {
		return fmt.Errorf("unable to create kubernetes secret: %w", err)
	}
	fmt.Printf("Applied secret %s\n", s.Name)

	podsClient := clientset.CoreV1().Pods(r.KubeNamespace)

	podName := util.GetPodName(chaincodeData.MspID, r.PeerID, chaincodeData.ChaincodeID)

	pod := util.GetChaincodePodObject(imageData, r.KubeNamespace, podName, r.PeerID, chaincodeData)

	createAttempts := 0
	for {
		createAttempts += 1
		p, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				if createAttempts > 3 {
					// give up
					return fmt.Errorf("unable to create chaincode pod %s/%s on final attempt: %w", r.KubeNamespace, podName, err)
				}

				err = podsClient.Delete(context.TODO(), podName, metav1.DeleteOptions{})
				if err != nil {
					if !errors.IsNotFound(err) {
						fmt.Fprintf(os.Stderr, "Error deleting existing chaincode pod: %v", err)
					}
				}

				_, err := util.WaitForPodTermination(context.TODO(), time.Minute, podsClient, podName, r.KubeNamespace)
				if err != nil {
					if !errors.IsNotFound(err) {
						fmt.Fprintf(os.Stderr, "Error waiting for existing chaincode pod to terminate: %v", err)
					}
				}

				// try again
				continue
			}

			return fmt.Errorf("unable to create chaincode pod %s/%s: %w", r.KubeNamespace, podName, err)
		}

		fmt.Printf("Created chaincode pod: %s/%s\n", p.Namespace, p.Name)
		break
	}

	_, err = util.WaitForPodRunning(context.TODO(), time.Minute, podsClient, podName, r.KubeNamespace)
	if err != nil {
		return fmt.Errorf("error waiting for chaincode pod: %w", err)
	}

	status, err := util.WaitForPodTermination(context.TODO(), 0, podsClient, podName, r.KubeNamespace)
	if err != nil {
		return fmt.Errorf("error waiting for chaincode pod to terminate: %w", err)
	}
	if status != nil {
		return fmt.Errorf("chaincode pod terminated %s: %s", status.Reason, status.Message)
	}

	return fmt.Errorf("unexpected chaincode pod termination")
}
