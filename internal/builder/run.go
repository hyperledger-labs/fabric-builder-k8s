// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/base64"
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

const (
	// Mutual TLS auth client key and cert paths in the chaincode container
	TLSClientKeyPath      string = "/etc/hyperledger/fabric/client.key"
	TLSClientCertPath     string = "/etc/hyperledger/fabric/client.crt"
	TLSClientKeyFile      string = "/etc/hyperledger/fabric/client_pem.key"
	TLSClientCertFile     string = "/etc/hyperledger/fabric/client_pem.crt"
	TLSClientRootCertFile string = "/etc/hyperledger/fabric/peer.crt"
)

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

	fmt.Fprintf(os.Stdout, "Chaincode ID: %s\n", chaincodeData.ChaincodeID)

	clientset, err := util.GetKubeClientset(r.KubeconfigPath)
	if err != nil {
		return fmt.Errorf("unable to connect kubernetes client: %v", err)
	}

	secretsClient := clientset.CoreV1().Secrets(r.KubeNamespace)

	// TODO need better/safer secret name!
	secretName := r.PeerID + "-secret-" + chaincodeData.ChaincodeID

	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
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

	_, err = secretsClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to create kubernetes secret: %v", err)
	}
	fmt.Printf("Created secret %s\n", secretName)

	deploymentsClient := clientset.AppsV1().Deployments(r.KubeNamespace)

	// TODO need better/safer application name! There are restrictions!
	applicationName := r.PeerID + "-cc-" + chaincodeData.ChaincodeID
	chaincodeImage := imageData.Name + ":" + imageData.Tag

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: applicationName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": applicationName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": applicationName,
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
							// TODO ports?!
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "certs",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: secretName,
									// Items: []apiv1.KeyToPath{
									// },
									// DefaultMode: 0400,
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
	_, err = deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("unable to create kubernetes deployment: %v", err)
	}
	fmt.Printf("Created deployment %s.\n", applicationName)

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
