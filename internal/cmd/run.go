// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/hyperledger-labs/fabric-builder-k8s/internal/builder"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/log"
	"github.com/hyperledger-labs/fabric-builder-k8s/internal/util"
	apiv1 "k8s.io/api/core/v1"
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
func getKubeHostAliases(logger *log.CmdLogger) (hostAliases []apiv1.HostAlias, ok bool) {
	raw := util.GetOptionalEnv(util.ChaincodeHostAliasesVariable, "")
	logger.Debugf("%s=%s", util.ChaincodeHostAliasesVariable, raw)

	if raw == "" {
		return nil, true
	}

	if err := json.Unmarshal([]byte(raw), &hostAliases); err != nil {
		logger.Printf(
			`The %s environment variable must be a valid JSON array, e.g. [{"ip":"1.2.3.4","hostnames":["foo.com"]}]: %v`,
			util.ChaincodeHostAliasesVariable, err,
		)

		return nil, false
	}

	// Validate IP addresses in host aliases
	for i, hostAlias := range hostAliases {
		if hostAlias.IP == "" {
			logger.Printf("The %s environment variable contains a host alias at index %d with an empty IP address", util.ChaincodeHostAliasesVariable, i)
			return nil, false
		}

		if !util.IsValidIPAddress(hostAlias.IP) {
			logger.Printf("The %s environment variable contains an invalid IP address '%s' at index %d", util.ChaincodeHostAliasesVariable, hostAlias.IP, i)
			return nil, false
		}

		if len(hostAlias.Hostnames) == 0 {
			logger.Printf("The %s environment variable contains a host alias at index %d with no hostnames", util.ChaincodeHostAliasesVariable, i)
			return nil, false
		}
	}

	return hostAliases, true
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

//nolint:nonamedreturns // using the ok bool convention to indicate errors
func getNameServers(logger *log.CmdLogger) (nameServers string, ok bool) {
	nameServers = util.GetOptionalEnv(util.NameServersVariable, "")
	logger.Debugf("%s=%s", util.NameServersVariable, nameServers)

	if nameServers == "" {
		return nameServers, true
	}

	// Validate IP address format
	if !util.IsValidIPAddress(nameServers) {
		logger.Printf("The %s environment variable must be a valid IP address", util.NameServersVariable)
		return "", false
	}

	return nameServers, true
}

//nolint:nonamedreturns // using the ok bool convention to indicate errors
func getCustomAnnotations(logger *log.CmdLogger) (annotations map[string]string, ok bool) {
	annotationsStr := util.GetOptionalEnv(util.CustomAnnotationsVariable, "")
	logger.Debugf("%s=%s", util.CustomAnnotationsVariable, annotationsStr)

	if annotationsStr == "" {
		return make(map[string]string), true
	}

	annotations = util.ParseAnnotations(annotationsStr)

	// Validate annotation keys follow Kubernetes naming conventions
	for key := range annotations {
		if !util.IsValidAnnotationKey(key) {
			logger.Printf("The %s environment variable contains an invalid annotation key '%s': must be a valid Kubernetes annotation key", util.CustomAnnotationsVariable, key)
			return nil, false
		}
	}

	logger.Debugf("Parsed custom annotations: %v", annotations)

	return annotations, true
}

func Run() {
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

		os.Exit(1)
	}

	buildOutputDirectory := os.Args[buildOutputDirectoryArg]
	runMetadataDirectory := os.Args[runMetadataDirectoryArg]

	logger.Debugf("Build output directory: %s", buildOutputDirectory)
	logger.Debugf("Run metadata directory: %s", runMetadataDirectory)

	//nolint:varnamelen // using the ok bool convention to indicate errors
	var ok bool

	peerID, ok := getPeerID(logger)
	if !ok {
		os.Exit(1)
	}

	kubeconfigPath := getKubeconfigPath(logger)
	kubeNamespace := getKubeNamespace(logger)

	kubeNodeRole, ok := getKubeNodeRole(logger)
	if !ok {
		os.Exit(1)
	}

	kubeServiceAccount := getKubeServiceAccount(logger)

	kubeNamePrefix, ok := getKubeNamePrefix(logger)
	if !ok {
		os.Exit(1)
	}

	chaincodeStartTimeout, ok := getChaincodeStartTimeout(logger)
	if !ok {
		os.Exit(1)
	}

	nameServers, ok := getNameServers(logger)
	if !ok {
		os.Exit(1)
	}

	customAnnotations, ok := getCustomAnnotations(logger)
	if !ok {
		os.Exit(1)
	}

	kubeHostAliases, ok := getKubeHostAliases(logger)
	if !ok {
		os.Exit(1)
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
		NameServers:           nameServers,
		CustomAnnotations:     customAnnotations,
		KubeHostAliases:       kubeHostAliases,
	}

	if err := run.Run(ctx); err != nil {
		logger.Printf("Error running chaincode: %+v", err)

		os.Exit(1)
	}

	os.Exit(0)
}
