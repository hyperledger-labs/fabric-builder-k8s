# Objectives

The aim is for the k8s builder to work as closely as possible with the existing [Fabric chaincode lifecycle](https://hyperledger-fabric.readthedocs.io/en/latest/chaincode_lifecycle.html), making sensible compromises for deploying chaincode on Kubernetes within those limitations.
(The assumption being that there are more people with Kubernetes skills than are familiar with the inner workings of Fabric!)

The two key principles are:

1. **The contents of the chaincode package must uniquely identify the chaincode functions executed on the ledger:**  
   In the case of the k8s builder the chaincode source code is not actually inside the package.
   In order not to break the Fabric chaincode lifecycle, the chaincode image must be specified using an immutable `@digest`, not `:label` which can be altered post commit.
   See [Pull an image by digest (immutable identifier)](https://docs.docker.com/engine/reference/commandline/pull/#pull-an-image-by-digest-immutable-identifier) for more details.

2. **The Fabric peer manages the chaincode process, not Kubernetes:** 
   Running the chaincode in server mode, i.e. allowing the peer to initiate the gRPC connection, would make it possible to leave Kubernetes to manage the chaincode process by creating a chaincode deployment.
   Unfortunately due to limitations in Fabric's builder and launcher implementation, that is not possible and the peer expects to control the chaincode process.

## Status

The k8s builder is [close to a version 1 release](https://github.com/hyperledger-labs/fabric-builder-k8s/milestone/1) and has been tested in a number of Kubernetes environments, deployment platforms, and provides semantic-revision aware [release tags](https://github.com/hyperledger-labs/fabric-builder-k8s/tags) for the external builder binaries.

The current status should be considered as STABLE and any bugs or enhancements delivered as GitHub Issues in conjunction with community PRs.
