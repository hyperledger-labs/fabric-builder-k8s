# Introduction

Kubernetes [external chaincode builder](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) for Hyperledger Fabric.

With the k8s builder, the Fabric administrator is responsible for [preparing a chaincode image](#chaincode-docker-image), publishing to a container registry, and [preparing a chaincode package](#chaincode-package) with coordinates of the contract's immutable image digest.
When Fabric detects the installation of a `type=k8s` contract, the builder assumes full ownership of the lifecycle of pods, containers, and network linkages necessary to communicate securely with the peer.


Advantages:

🚀 Chaincode runs _immediately_ on channel commit.

✨ Avoids the complexity and administrative burdens associated with Chaincode-as-a-Service.

🔥 Pre-published chaincode images avoid code-compilation errors at deployment time.

🏗️ Pre-published chaincode images encourage modern, industry accepted CI/CD best practices.

🛡️ Pre-published chaincode images remove any and all dependencies on a root-level _docker daemon_.

🕵️ Pre-published chaincode images provide traceability and change management features (e.g. Git commit hash as image tag)

The aim is for the builder to work as closely as possible with the existing [Fabric chaincode lifecycle](https://hyperledger-fabric.readthedocs.io/en/latest/chaincode_lifecycle.html), making sensible compromises for deploying chaincode on Kubernetes within those constraints.
(The assumption being that there are more people with Kubernetes skills than are familiar with the inner workings of Fabric!)

For example:

- The contents of the chaincode package must uniquely identify the chaincode functions executed on the ledger. 

  In the case of the k8s builder the chaincode source code is not actually inside the package.  In order not to break the Fabric chaincode lifecycle, the chaincode image must be specified using an immutable `@digest`, not a `:label` which can be altered post commit.
  
  See [Pull an image by digest (immutable identifier)](https://docs.docker.com/engine/reference/commandline/pull/#pull-an-image-by-digest-immutable-identifier) for more details.


- The Fabric peer manages the chaincode process, not Kubernetes.

  Running the chaincode in server mode, i.e. allowing the peer to initiate the gRPC connection, would make it possible to leave Kubernetes to manage the chaincode process by creating a chaincode deployment.

  Unfortunately due to limitations in Fabric's builder and launcher implementation, that is not possible and the peer expects to control the chaincode process.


**Status:** the k8s builder is [close to a version 1 release](https://github.com/hyperledger-labs/fabric-builder-k8s/milestone/1) and has been tested in a number of Kubernetes environments, deployment platforms, and provides semantic-revision aware [release tags](https://github.com/hyperledger-labs/fabric-builder-k8s/tags) for the external builder binaries.
The current status should be considered as STABLE and any bugs or enhancements delivered as GitHub Issues in conjunction with community PRs.



There are addition docs with more detailed usage instructions for specific Fabric network deployments:

- [Kubernetes Test Network](tutorials/test-network-k8s.md)
- [Nano Test Network](tutorials/test-network-nano.md)
- [Fabric Operator](tutorials/fabric-operator.md)
- [HLF Operator](tutorials/hlf-operator.md)




## Chaincode deploy

Deploy the chaincode package as usual, starting by installing the k8s chaincode package.

```shell
peer lifecycle chaincode install go-contract.tgz
```

You can also use the `peer` command to get the chaincode package ID.

```shell
export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid go-contract.tgz) && echo $PACKAGE_ID
```
