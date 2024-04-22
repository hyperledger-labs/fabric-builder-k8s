// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
)

type Run struct {
	BuildOutputDirectory string
	RunMetadataDirectory string
	PeerID               string
	KubeconfigPath       string
	KubeNamespace        string
	KubeServiceAccount   string
	KubeNamePrefix       string
}

func (r *Run) Run(ctx context.Context) error {
	logger := log.New(ctx)
	logger.Debugln("Running chaincode...")

	imageData, err := util.ReadImageJSON(logger, r.BuildOutputDirectory)
	if err != nil {
		return err
	}

	chaincodeData, err := util.ReadChaincodeJSON(logger, r.RunMetadataDirectory)
	if err != nil {
		return err
	}

	kubeObjectName := util.GetValidRfc1035LabelName(r.KubeNamePrefix, r.PeerID, chaincodeData)

	clientset, err := util.GetKubeClientset(logger, r.KubeconfigPath)
	if err != nil {
		return fmt.Errorf(
			"unable to connect kubernetes client for chaincode ID %s: %w",
			chaincodeData.ChaincodeID,
			err,
		)
	}

	secretsClient := clientset.CoreV1().Secrets(r.KubeNamespace)

	err = util.ApplyChaincodeSecrets(
		ctx,
		logger,
		secretsClient,
		kubeObjectName,
		r.KubeNamespace,
		r.PeerID,
		chaincodeData,
	)
	if err != nil {
		return fmt.Errorf(
			"unable to create kubernetes secret for chaincode ID %s: %w",
			chaincodeData.ChaincodeID,
			err,
		)
	}

	podsClient := clientset.CoreV1().Pods(r.KubeNamespace)

	pod, err := util.CreateChaincodePod(
		ctx,
		logger,
		podsClient,
		kubeObjectName,
		r.KubeNamespace,
		r.KubeServiceAccount,
		r.PeerID,
		chaincodeData,
		imageData,
	)
	if err != nil {
		return err
	}

	logger.Printf(
		"Running chaincode ID %s in kubernetes pod %s/%s",
		chaincodeData.ChaincodeID,
		pod.Namespace,
		pod.Name,
	)

	return util.WaitForChaincodePod(ctx, logger, podsClient, pod, chaincodeData.ChaincodeID)
}
