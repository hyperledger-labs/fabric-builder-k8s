// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"os"
	"time"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/builder"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation"
)

//nolint:nonamedreturns // using the ok bool convention to indicate errors
func getPeerID(logger *log.CmdLogger) (peerID string, ok bool) {
	peerID, err := util.GetRequiredEnv(util.PeerIDVariable)
	if err != nil {
		logger.Printf("Expected %s environment variable\n", util.PeerIDVariable)

		return peerID, false
	}

	logger.Debugf("%s=%s", util.PeerIDVariable, peerID)

	return peerID, true
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

//nolint:nonamedreturns // using the ok bool convention to indicate errors
func getKubeNodeRole(logger *log.CmdLogger) (kubeNodeRole string, ok bool) {
	kubeNodeRole = util.GetOptionalEnv(util.ChaincodeNodeRoleVariable, "")
	logger.Debugf("%s=%s", util.ChaincodeNodeRoleVariable, kubeNodeRole)

	// TODO: are valid taint values the same?!
	if msgs := validation.IsValidLabelValue(kubeNodeRole); len(msgs) > 0 {
		logger.Printf("The %s environment variable must be a valid Kubernetes label value: %s", util.ChaincodeNodeRoleVariable, msgs[0])

		return kubeNodeRole, false
	}

	return kubeNodeRole, true
}

func getKubeServiceAccount(logger *log.CmdLogger) string {
	kubeServiceAccount := util.GetOptionalEnv(util.ChaincodeServiceAccountVariable, util.DefaultServiceAccountName)
	logger.Debugf("%s=%s", util.ChaincodeServiceAccountVariable, kubeServiceAccount)

	return kubeServiceAccount
}

//nolint:nonamedreturns // using the ok bool convention to indicate errors
func getKubeNamePrefix(logger *log.CmdLogger) (kubeNamePrefix string, ok bool) {
	const maximumKubeNamePrefixLength = 30

	kubeNamePrefix = util.GetOptionalEnv(util.ObjectNamePrefixVariable, util.DefaultObjectNamePrefix)
	logger.Debugf("%s=%s", util.ObjectNamePrefixVariable, kubeNamePrefix)

	if len(kubeNamePrefix) > maximumKubeNamePrefixLength {
		logger.Printf("The %s environment variable must be a maximum of 30 characters", util.ObjectNamePrefixVariable)

		return kubeNamePrefix, false
	}

	if msgs := apivalidation.NameIsDNS1035Label(kubeNamePrefix, true); len(msgs) > 0 {
		logger.Printf("The %s environment variable must be a valid DNS-1035 label: %s", util.ObjectNamePrefixVariable, msgs[0])

		return kubeNamePrefix, false
	}

	return kubeNamePrefix, true
}

//nolint:nonamedreturns // using the ok bool convention to indicate errors
func getChaincodeStartTimeout(logger *log.CmdLogger) (chaincodeStartTimeoutDuration time.Duration, ok bool) {
	chaincodeStartTimeout := util.GetOptionalEnv(util.ChaincodeStartTimeoutVariable, util.DefaultStartTimeout)
	logger.Debugf("%s=%s", util.ChaincodeStartTimeoutVariable, chaincodeStartTimeout)

	chaincodeStartTimeoutDuration, err := time.ParseDuration(chaincodeStartTimeout)
	if err != nil {
		logger.Printf("The %s environment variable must be a valid Go duration string, e.g. 3m40s: %v", util.ChaincodeStartTimeoutVariable, err)

		return 0 * time.Minute, false
	}

	return chaincodeStartTimeoutDuration, true
}

func Run() int {
	const (
		expectedArgsLength      = 3
		buildOutputDirectoryArg = 1
		runMetadataDirectoryArg = 2
	)

	debug := util.GetOptionalEnv(util.DebugVariable, "false")
	ctx := log.NewCmdContext(context.Background(), debug == "true")
	logger := log.New(ctx)

	if len(os.Args) != expectedArgsLength {
		logger.Println("Expected BUILD_OUTPUT_DIR and RUN_METADATA_DIR arguments")

		return 1
	}

	buildOutputDirectory := os.Args[buildOutputDirectoryArg]
	runMetadataDirectory := os.Args[runMetadataDirectoryArg]

	logger.Debugf("Build output directory: %s", buildOutputDirectory)
	logger.Debugf("Run metadata directory: %s", runMetadataDirectory)

	//nolint:varnamelen // using the ok bool convention to indicate errors
	var ok bool

	peerID, ok := getPeerID(logger)
	if !ok {
		return 1
	}

	kubeconfigPath := getKubeconfigPath(logger)
	kubeNamespace := getKubeNamespace(logger)

	kubeNodeRole, ok := getKubeNodeRole(logger)
	if !ok {
		return 1
	}

	kubeServiceAccount := getKubeServiceAccount(logger)

	kubeNamePrefix, ok := getKubeNamePrefix(logger)
	if !ok {
		return 1
	}

	chaincodeStartTimeout, ok := getChaincodeStartTimeout(logger)
	if !ok {
		return 1
	}

	run := &builder.Run{
		BuildOutputDirectory:  buildOutputDirectory,
		RunMetadataDirectory:  runMetadataDirectory,
		PeerID:                peerID,
		KubeconfigPath:        kubeconfigPath,
		KubeNamespace:         kubeNamespace,
		KubeNodeRole:          kubeNodeRole,
		KubeServiceAccount:    kubeServiceAccount,
		KubeNamePrefix:        kubeNamePrefix,
		ChaincodeStartTimeout: chaincodeStartTimeout,
	}

	if err := run.Run(ctx); err != nil {
		logger.Printf("Error running chaincode: %+v", err)

		return 1
	}

	return 0
}
