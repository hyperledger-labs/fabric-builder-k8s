// SPDX-License-Identifier: Apache-2.0

package util

import (
	"fmt"
	"net"
	"os"
	"regexp"
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
	ChaincodeHostAliasesVariable    = builderVariablePrefix + "HOST_ALIASES"
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

// IsValidIPAddress validates if a string is a valid IPv4 or IPv6 address.
func IsValidIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsValidAnnotationKey validates if a string is a valid Kubernetes annotation key.
// Annotation keys must be in the format: [prefix/]name
// - prefix (optional): DNS subdomain (max 253 chars)
// - name (required): max 63 chars, alphanumeric, '-', '_', '.'
func IsValidAnnotationKey(key string) bool {
	if key == "" {
		return false
	}

	// Split into prefix and name
	parts := strings.SplitN(key, "/", 2)
	
	var prefix, name string
	if len(parts) == 2 {
		prefix = parts[0]
		name = parts[1]
	} else {
		name = parts[0]
	}

	// Validate prefix if present (DNS subdomain format)
	if prefix != "" {
		if len(prefix) > 253 {
			return false
		}
		// DNS subdomain regex: lowercase alphanumeric, '-', '.'
		dnsSubdomainRegex := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
		if !dnsSubdomainRegex.MatchString(prefix) {
			return false
		}
	}

	// Validate name (required)
	if name == "" || len(name) > 63 {
		return false
	}

	// Name must be alphanumeric, '-', '_', '.' and start/end with alphanumeric
	nameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([-a-zA-Z0-9_.]*[a-zA-Z0-9])?$`)
	return nameRegex.MatchString(name)
}
