#!/usr/bin/env sh
#
# SPDX-License-Identifier: Apache-2.0
#
set -eu

mkdir -p "${HOME}"/.local/bin

export GOENV_OS=$(go env GOOS)
export GOENV_ARCH=$(go env GOARCH)
export UNAME_KERNAL=$(uname -s)

#
# Install k9s
#
curl -sSL https://github.com/derailed/k9s/releases/download/v0.32.4/k9s_${UNAME_KERNAL}_${GOENV_ARCH}.tar.gz | tar -zxf - -C "${HOME}/.local/bin/" k9s && chmod +x "${HOME}/.local/bin/k9s"

#
# Install yq
#
curl -sSLo "${HOME}/.local/bin/yq" https://github.com/mikefarah/yq/releases/download/v4.43.1/yq_${GOENV_OS}_${GOENV_ARCH} && chmod +x "${HOME}/.local/bin/yq"

#
# Install fabric binaries and the nano test network
#
rm -r "${PWD}"/.fabric || true
mkdir "${PWD}"/.fabric
cd .fabric

curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh binary

export FABRIC_SAMPLES_COMMIT=0db64487e5e89a81d68e6871af3f0907c67e7d75
curl -sSL "https://github.com/hyperledger/fabric-samples/archive/${FABRIC_SAMPLES_COMMIT}.tar.gz" | tar -xzf - --strip-components=1 fabric-samples-${FABRIC_SAMPLES_COMMIT}/test-network-nano-bash

cd ..

#
# Add k8s builder config to fabric core.yaml
#
# To install the k8s builder use the following command:
#   GOBIN="${PWD}"/.fabric/builders/k8s_builder/bin go install ./cmd/...
#
export FABRIC_K8S_BUILDER_PATH="${PWD}/.fabric/builders/k8s_builder"
mkdir -p "${FABRIC_K8S_BUILDER_PATH}/bin"

yq -i '.chaincode.externalBuilders += { "name": "k8s_builder", "path": "${FABRIC_K8S_BUILDER_PATH}" | envsubst(ne), "propagateEnvironment": [ "CORE_PEER_ID", "KUBECONFIG_PATH", "FABRIC_K8S_BUILDER_DEBUG" ] }' .fabric/config/core.yaml
