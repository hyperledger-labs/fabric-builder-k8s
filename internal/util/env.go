// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	builderVariablePrefix           = "FABRIC_K8S_BUILDER_"
	ChaincodeNamespaceVariable      = builderVariablePrefix + "NAMESPACE"
	ChaincodeNodeRoleVariable       = builderVariablePrefix + "NODE_ROLE"
	ObjectNamePrefixVariable        = builderVariablePrefix + "OBJECT_NAME_PREFIX"
	ChaincodeServiceAccountVariable = builderVariablePrefix + "SERVICE_ACCOUNT"
	ChaincodeStartTimeoutVariable   = builderVariablePrefix + "START_TIMEOUT"
	NameServersVariable             = builderVariablePrefix + "NAME_SERVERS"
	CustomAnnotationsVariable       = builderVariablePrefix + "CUSTOM_ANNOTATIONS"
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

// ParseAnnotations parses a comma-separated list of key=value pairs into a map.
// Example input: "sidecar.istio.io/inject=true,app=myapp"
// Returns empty map if input is empty or invalid entries are skipped.
func ParseAnnotations(annotationsStr string) map[string]string {
	annotations := make(map[string]string)

	if annotationsStr == "" {
		return annotations
	}

	pairs := strings.Split(annotationsStr, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				annotations[key] = value
			}
		}
	}

	return annotations
}
