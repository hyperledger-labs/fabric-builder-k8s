// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/builder"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
)

const (
	expectedArgsLength      = 3
	buildOutputDirectoryArg = 1
	runMetadataDirectoryArg = 2
)

func main() {
	debug := util.GetOptionalEnv(util.DebugVariable, "false")
	ctx := log.NewCmdContext(context.Background(), debug == "true")
	logger := log.New(ctx)

	if len(os.Args) != expectedArgsLength {
		logger.Println("Expected BUILD_OUTPUT_DIR and RUN_METADATA_DIR arguments")
		os.Exit(1)
	}

	buildOutputDirectory := os.Args[buildOutputDirectoryArg]
	runMetadataDirectory := os.Args[runMetadataDirectoryArg]

	logger.Debugf("Build output directory: %s", buildOutputDirectory)
	logger.Debugf("Run metadata directory: %s", runMetadataDirectory)

	peerID, err := util.GetRequiredEnv(util.PeerIDVariable)
	if err != nil {
		logger.Printf("Expected %s environment variable\n", util.PeerIDVariable)
		os.Exit(1)
	}

	logger.Debugf("%s=%s", util.PeerIDVariable, peerID)

	kubeconfigPath := util.GetOptionalEnv(util.KubeconfigPathVariable, "")
	logger.Debugf("%s=%s", util.KubeconfigPathVariable, kubeconfigPath)

	kubeNamespace := util.GetOptionalEnv(util.ChaincodeNamespaceVariable, "")
	logger.Debugf("%s=%s", util.ChaincodeNamespaceVariable, kubeNamespace)

	if kubeNamespace == "" {
		kubeNamespace, err = util.GetKubeNamespace()
		if err != nil {
			kubeNamespace = "default"
		}
	}

	kubeServiceAccount := util.GetOptionalEnv(util.ChaincodeServiceAccountVariable, "default")
	logger.Debugf("%s=%s", util.ChaincodeServiceAccountVariable, kubeServiceAccount)

	run := &builder.Run{
		BuildOutputDirectory: buildOutputDirectory,
		RunMetadataDirectory: runMetadataDirectory,
		PeerID:               peerID,
		KubeconfigPath:       kubeconfigPath,
		KubeNamespace:        kubeNamespace,
		KubeServiceAccount:   kubeServiceAccount,
	}

	if err := run.Run(ctx); err != nil {
		logger.Printf("Error running chaincode: %+v", err)
		os.Exit(1)
	}
}
