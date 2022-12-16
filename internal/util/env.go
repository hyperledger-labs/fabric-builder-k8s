// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"os"
)

const (
	builderVariablePrefix           = "FABRIC_K8S_BUILDER_"
	ChaincodeNamespaceVariable      = builderVariablePrefix + "NAMESPACE"
	ChaincodeServiceAccountVariable = builderVariablePrefix + "SERVICE_ACCOUNT"
	DebugVariable                   = builderVariablePrefix + "DEBUG"
	KubeconfigPathVariable          = "KUBECONFIG_PATH"
	PeerIDVariable                  = "CORE_PEER_ID"
)

func GetOptionalEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func GetRequiredEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}

	return "", fmt.Errorf("environment variable not set: %s", key)
}
