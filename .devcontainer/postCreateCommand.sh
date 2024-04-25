#!/usr/bin/env sh
#
# SPDX-License-Identifier: Apache-2.0
#
set -eu

export GOOS=$(go env GOOS)
export GOARCH=$(go env GOARCH)

#
# Install yq
#
mkdir -p "${HOME}"/.local/bin
curl -sSLo "${HOME}"/.local/bin/yq https://github.com/mikefarah/yq/releases/download/v4.43.1/yq_${GOOS}_${GOARCH} && chmod +x "${HOME}"/.local/bin/yq

#
# Install fabric binaries and the nano test network
#
rm -r "${PWD}"/.fabric || true
mkdir "${PWD}"/.fabric
cd .fabric

curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh binary

export FABRIC_SAMPLES_COMMIT=20009ecd029fa887a8a98122fa5df0ec5181cdb1
curl -sSL "https://github.com/hyperledger/fabric-samples/archive/${FABRIC_SAMPLES_COMMIT}.tar.gz" | tar -xzf - --strip-components=1 fabric-samples-${FABRIC_SAMPLES_COMMIT}/test-network-nano-bash

cd ..

#
# Add k8s builder config to fabric core.yaml
#
# To install the k8s builder use the following command:
#   GOBIN="${PWD}"/.fabric/builders/k8s_builder go install ./cmd/...
#
export FABRIC_K8S_BUILDER_PATH="${PWD}"/.fabric/builders/k8s_builder
mkdir -p "${FABRIC_K8S_BUILDER_PATH}"

yq -i '.chaincode.externalBuilders += { "name": "k8s_builder", "path": "${FABRIC_K8S_BUILDER_PATH}" | envsubst(ne), "propagateEnvironment": [ "CORE_PEER_ID", "KUBECONFIG_PATH" ] }' .fabric/config/core.yaml
