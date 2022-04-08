// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
}

func (r *Run) Run() error {
	imageJsonPath := filepath.Join(r.BuildOutputDirectory, "/image.json")
	chaincodeJsonPath := filepath.Join(r.RunMetadataDirectory, "/chaincode.json")

	// Read the image.json file
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

	// Read the chaincode.json file
	fmt.Println("Reading chaincode.json...")
	_, err = os.Stat(chaincodeJsonPath)
	if err != nil {
		return err
	}

	chaincodeJsonContents, err := ioutil.ReadFile(chaincodeJsonPath)
	if err != nil {
		return err
	}

	var chaincodeData chaincode
	if err := json.Unmarshal(chaincodeJsonContents, &chaincodeData); err != nil {
		fmt.Errorf(err.Error())
		return err
	}

	fmt.Fprintf(os.Stdout, "Chaincode ID: %s", chaincodeData.ChaincodeID)

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// TODO use namespace from env var?
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	// TODO read in a deployment YAML if specified?
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			// TODO unique name!
			Name: "chaincode-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			// TODO labels?
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			// TODO??
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: "chaincode",
							// TODO use imageData!
							Image: "nginx:1.12",
							// TODO ports, env, etc.
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
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
