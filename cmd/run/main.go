// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/builder"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation"
)

const (
	expectedArgsLength          = 3
	buildOutputDirectoryArg     = 1
	runMetadataDirectoryArg     = 2
	maximumKubeNamePrefixLength = 30
)

func getPeerID(logger *log.CmdLogger) string {
	peerID, err := util.GetRequiredEnv(util.PeerIDVariable)
	if err != nil {
		logger.Printf("Expected %s environment variable\n", util.PeerIDVariable)
		os.Exit(1)
	}

	logger.Debugf("%s=%s", util.PeerIDVariable, peerID)

	return peerID
}

func getKubeconfigPath(logger *log.CmdLogger) string {
	kubeconfigPath := util.GetOptionalEnv(util.KubeconfigPathVariable, "")
	logger.Debugf("%s=%s", util.KubeconfigPathVariable, kubeconfigPath)

	return kubeconfigPath
}

func getKubeNamespace(logger *log.CmdLogger) string {
	kubeNamespace := util.GetOptionalEnv(util.ChaincodeNamespaceVariable, "")
	logger.Debugf("%s=%s", util.ChaincodeNamespaceVariable, kubeNamespace)

	if kubeNamespace == "" {
		var err error

		kubeNamespace, err = util.GetKubeNamespace()
		if err != nil {
			logger.Debugf("Error getting namespace: %+v\n", util.DefaultNamespace, err)
			kubeNamespace = util.DefaultNamespace
		}

		logger.Debugf("Using default namespace: %s\n", util.DefaultNamespace)
	}

	return kubeNamespace
}

func getKubeNodeRole(logger *log.CmdLogger) string {
	kubeNodeRole := util.GetOptionalEnv(util.ChaincodeNodeRoleVariable, "")
	logger.Debugf("%s=%s", util.ChaincodeNodeRoleVariable, kubeNodeRole)

	// TODO: are valid taint values the same?!
	if msgs := validation.IsValidLabelValue(kubeNodeRole); len(msgs) > 0 {
		logger.Printf("The %s environment variable must be a valid Kubernetes label value: %s", util.ChaincodeNodeRoleVariable, msgs[0])
		os.Exit(1)
	}

	return kubeNodeRole
}

func getKubeServiceAccount(logger *log.CmdLogger) string {
	kubeServiceAccount := util.GetOptionalEnv(util.ChaincodeServiceAccountVariable, util.DefaultServiceAccountName)
	logger.Debugf("%s=%s", util.ChaincodeServiceAccountVariable, kubeServiceAccount)

	return kubeServiceAccount
}

func getKubeNamePrefix(logger *log.CmdLogger) string {
	kubeNamePrefix := util.GetOptionalEnv(util.ObjectNamePrefixVariable, util.DefaultObjectNamePrefix)
	logger.Debugf("%s=%s", util.ObjectNamePrefixVariable, kubeNamePrefix)

	if len(kubeNamePrefix) > maximumKubeNamePrefixLength {
		logger.Printf("The %s environment variable must be a maximum of 30 characters", util.ObjectNamePrefixVariable)
		os.Exit(1)
	}

	if msgs := apivalidation.NameIsDNS1035Label(kubeNamePrefix, true); len(msgs) > 0 {
		logger.Printf("The %s environment variable must be a valid DNS-1035 label: %s", util.ObjectNamePrefixVariable, msgs[0])
		os.Exit(1)
	}

	return kubeNamePrefix
}

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

	peerID := getPeerID(logger)
	kubeconfigPath := getKubeconfigPath(logger)
	kubeNamespace := getKubeNamespace(logger)
	kubeNodeRole := getKubeNodeRole(logger)
	kubeServiceAccount := getKubeServiceAccount(logger)
	kubeNamePrefix := getKubeNamePrefix(logger)

	run := &builder.Run{
		BuildOutputDirectory: buildOutputDirectory,
		RunMetadataDirectory: runMetadataDirectory,
		PeerID:               peerID,
		KubeconfigPath:       kubeconfigPath,
		KubeNamespace:        kubeNamespace,
		KubeNodeRole:         kubeNodeRole,
		KubeServiceAccount:   kubeServiceAccount,
		KubeNamePrefix:       kubeNamePrefix,
	}

	if err := run.Run(ctx); err != nil {
		logger.Printf("Error running chaincode: %+v", err)
		os.Exit(1)
	}
}
