// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

// TODO: read these from the input scope / json / CLI params / env / etc.
const (
	Namespace = "test-network"
	PeerName = "peer1"
	ChaincodeName = "asset-transfer-basic"
	ChaincodeServerAddress = "0.0.0.0:9999"

	// todo: generate the "ApplicationName" for the chaincode endpoint
	ApplicationName = "org1" + PeerName + "-cc-" + ChaincodeName

	// todo: run locally or run in the cluster with the same client spec
	RunInCluster = true
)

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

	fmt.Fprintf(os.Stdout, "Image ApplicationName: %s\nImage tag: %s\n", imageData.Name, imageData.Tag)

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

	fmt.Printf("Chaincode ID: %s\n", chaincodeData.ChaincodeID)

	// todo: find a clientset constructor that will read from either the local .kube/config or the in-cluster service account
	var config *rest.Config

	if RunInCluster {
		// creates an in-cluster kube config
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}

	} else {
		// creates an out of cluster kube config
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// TODO use namespace from env var instead of "test-network"
	deploymentsClient := clientset.AppsV1().Deployments(Namespace)

	// TODO read in a deployment YAML if specified?
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: ApplicationName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
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
							Name: "main",
							Image: imageData.Name + ":" + imageData.Tag,
							Env: []apiv1.EnvVar{
								{
									Name: "CHAINCODE_SERVER_ADDRESS",
									Value: ChaincodeServerAddress,
								},
								{
									Name: "CHAINCODE_ID",
									Value: chaincodeData.ChaincodeID,
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "chaincode",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 9999,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	deploymentYaml, _ := yaml.Marshal(deployment)
	fmt.Printf("Creating deployment %s:\n", ApplicationName)
	fmt.Printf("%s\n", deploymentYaml)

	deployment, err = deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created deployment %q.\n", deployment.GetObjectMeta().GetName())


	// Create a Service to expose the chaincode endpoint to the peer
	servicesClient := clientset.CoreV1().Services(Namespace)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: ApplicationName,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": ApplicationName,
			},
			Ports: []apiv1.ServicePort{
				{
					Name: "chaincode",
					Port: 9999,
					Protocol: apiv1.ProtocolTCP,
				},
			},
		},
	}

	serviceYaml, _ := yaml.Marshal(service)
	fmt.Printf("Creating service %s:\n", ApplicationName)
	fmt.Printf("%s\n", serviceYaml)

	service, err = servicesClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Created service %q.\n", service.GetObjectMeta().GetName())

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
