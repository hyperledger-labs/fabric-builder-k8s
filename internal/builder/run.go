// SPDX-License-Identifier: Apache-2.0

package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
)

type Run struct {
	BuildOutputDirectory  string
	RunMetadataDirectory  string
	PeerID                string
	KubeconfigPath        string
	KubeNamespace         string
	KubeNodeRole          string
	KubeServiceAccount    string
	KubeNamePrefix        string
	ChaincodeStartTimeout time.Duration
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

	kubeObjectName := util.GetValidRfc1035LabelName(r.KubeNamePrefix, r.PeerID, chaincodeData, util.ObjectNameSuffixLength+1)

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

	jobsClient := clientset.BatchV1().Jobs(r.KubeNamespace)

	job, err := util.CreateChaincodeJob(
		ctx,
		logger,
		jobsClient,
		kubeObjectName,
		r.KubeNamespace,
		r.KubeServiceAccount,
		r.KubeNodeRole,
		r.PeerID,
		chaincodeData,
		imageData,
	)
	if err != nil {
		return err
	}

	logger.Printf(
		"Running chaincode ID %s with kubernetes job %s/%s",
		chaincodeData.ChaincodeID,
		job.Namespace,
		job.Name,
	)

	batchClient := clientset.BatchV1().RESTClient()

	return util.WaitForChaincodeJob(ctx, logger, batchClient, job, chaincodeData.ChaincodeID, r.ChaincodeStartTimeout)
}
