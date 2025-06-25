# Installation

Installing the k8s builder is part of the process of creating a production Fabric peer.

For more information, see the [Checklist for a production peer](https://hyperledger-fabric.readthedocs.io/en/latest/deploypeer/peerchecklist.html#chaincode-externalbuilders) Fabric documentation.

## Sample peer image

A sample [k8s-fabric-peer image](https://github.com/hyperledger-labs/fabric-builder-k8s/pkgs/container/fabric-builder-k8s%2Fk8s-fabric-peer) is available. The `k8s-fabric-peer` is based on the [hyperledger/fabric-peer](https://hub.docker.com/r/hyperledger/fabric-peer) with the k8s builder preconfigured.

## Prebuilt binaries

Prebuilt binaries are available to download from the [releases page](https://github.com/hyperledger-labs/fabric-builder-k8s/releases).

## Install from source

To install from source in `/opt/hyperledger/k8s_builder`, use the `go install` command.

```shell
mkdir -p /opt/hyperledger/k8s_builder/bin
cd /opt/hyperledger/k8s_builder/bin
GOBIN="${PWD}" go install github.com/hyperledger-labs/fabric-builder-k8s/cmd/...@v0.13.0
```
