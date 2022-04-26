// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hyperledgendary/fabric-builder-k8s/internal/util"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// chaincode represents the chaincode.json file that is supplied by Fabric in
// the RUN_METADATA_DIR
type chaincode struct {
	ChaincodeID string `json:"chaincode_id"`
	PeerAddress string `json:"peer_address"`
	ClientCert  string `json:"client_cert"`
	ClientKey   string `json:"client_key"`
	RootCert    string `json:"root_cert"`
	MspID       string `json:"mspid"`
}

// image represents the image.json file from the k8s chaincode tarball
type image struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
}

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
		return err
	}

	imageJsonContents, err := ioutil.ReadFile(imageJsonPath)
	if err != nil {
		return err
	}

	var imageData image
	if err := json.Unmarshal(imageJsonContents, &imageData); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Image name: %s\nImage tag: %s\n", imageData.Name, imageData.Tag)

	fmt.Println("Reading chaincode.json...")
	_, err = os.Stat(chaincodeJsonPath)
	if err != nil {
		return fmt.Errorf("unable to access chaincode.json: %v", err)
	}

	chaincodeJsonContents, err := ioutil.ReadFile(chaincodeJsonPath)
	if err != nil {
		return fmt.Errorf("unable to read chaincode.json: %v", err)
	}

	var chaincodeData chaincode
	if err := json.Unmarshal(chaincodeJsonContents, &chaincodeData); err != nil {
		return fmt.Errorf("unable to process chaincode.json: %v", err)
	}

	fmt.Fprintf(os.Stdout, "Chaincode ID: %s", chaincodeData.ChaincodeID)

	clientset, err := util.GetKubeClientset(r.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("unable to connect kubernetes client: %v", err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(r.KubeNamespace)

	ApplicationName := r.PeerID + "-cc-" + chaincodeData.ChaincodeID
	ChaincodeImage := imageData.Name + ":" + imageData.Tag

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: ApplicationName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": ApplicationName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": ApplicationName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "main",
							Image: ChaincodeImage,
							// TODO ports, env, etc.
							Env: []apiv1.EnvVar{
								{
									Name:  "CORE_CHAINCODE_ID_NAME",
									Value: chaincodeData.ChaincodeID,
								},
								{
									Name:  "CORE_PEER_ADDRESS",
									Value: chaincodeData.PeerAddress,
								},
								// {
								// 	Name: "CORE_PEER_TLS_ENABLED",
								// 	Value: "true",
								// },
								// {
								// 	Name: "CORE_PEER_TLS_ROOTCERT_FILE",
								// 	Value: "/certs/peer.crt",
								// },
								// {
								// 	Name: "CORE_TLS_CLIENT_KEY_PATH",
								// 	Value: "/certs/client.key",
								// },
								// {
								// 	Name: "CORE_TLS_CLIENT_CERT_PATH",
								// 	Value: "/certs/client.crt",
								// },
								// {
								// 	Name: "CORE_TLS_CLIENT_KEY_FILE",
								// 	Value: "/certs/client_pem.key",
								// },
								// {
								// 	Name: "CORE_TLS_CLIENT_CERT_FILE",
								// 	Value: "/certs/client_pem.crt",
								// },
								// {
								// 	Name: "CORE_PEER_LOCALMSPID",
								// 	Value: os.Getenv("CORE_PEER_LOCALMSPID")
								// },
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to create kubernetes deployment: %v", err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
