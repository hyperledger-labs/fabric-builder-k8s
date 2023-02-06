# fabric-builder-k8s

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

## Usage

The k8s builder can be run in cluster using the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` environment variables, or it can connect using a `KUBECONFIG_PATH` environment variable.

The following optional environment variables can be used to configure the k8s builder:

- `FABRIC_K8S_BUILDER_DEBUG` whether to enable additional logging
- `FABRIC_K8S_BUILDER_NAMESPACE` specifies the namespace to deploy chaincode to
- `FABRIC_K8S_BUILDER_SERVICE_ACCOUNT` specifies the service account for the chaincode pod

A `CORE_PEER_ID` environment variable is also currently required.

External builders are configured in the `core.yaml` file, for example:

```
  externalBuilders:
    - name: k8s_builder
      path: /opt/hyperledger/k8s_builder
      propagateEnvironment:
        - CORE_PEER_ID
        - FABRIC_K8S_BUILDER_DEBUG
        - FABRIC_K8S_BUILDER_NAMESPACE
        - FABRIC_K8S_BUILDER_SERVICE_ACCOUNT
        - KUBERNETES_SERVICE_HOST
        - KUBERNETES_SERVICE_PORT
```

See [External Builders and Launchers](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) for details of Hyperledger Fabric builders.

As well as configuring Fabric to use the k8s builder, you will need to [configure Kubernetes](docs/KUBERNETES_CONFIG.md) to allow the builder to start chaincode pods successfully.

There are addition docs with more detailed usage instructions for specific Fabric network deployments:

- [Kubernetes Test Network](docs/TEST_NETWORK_K8S.md)
- [Nano Test Network](docs/TEST_NETWORK_NANO.md)
- [Fabric Operator](docs/FABRIC_OPERATOR.md)
- [HLF Operator](docs/HLF_OPERATOR.md)

## Chaincode Docker image

Unlike the traditional chaincode language support for Go, Java, and Node.js, the k8s builder *does not* build a chaincode Docker image using Docker-in-Docker.
Instead, a chaincode Docker image must be built and published before it can be used with the k8s builder.

The chaincode will have access to the following environment variables:

- CORE_CHAINCODE_ID_NAME
- CORE_PEER_ADDRESS
- CORE_PEER_TLS_ENABLED
- CORE_PEER_TLS_ROOTCERT_FILE
- CORE_TLS_CLIENT_KEY_PATH
- CORE_TLS_CLIENT_CERT_PATH
- CORE_TLS_CLIENT_KEY_FILE
- CORE_TLS_CLIENT_CERT_FILE
- CORE_PEER_LOCALMSPID

See the [sample contracts for Go, Java, and Node.js](samples/README.md) for basic docker images which will work with the k8s builder.

## Chaincode package

The k8s chaincode package file, which is installed by the `peer lifecycle chaincode install` command, must contain the Docker image name and digest of the chaincode being deployed.

[Fabric chaincode packages](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html#chaincode-packages) are `.tgz` files which contain two files:

- metadata.json - the chaincode label and type
- code.tar.gz - source artifacts for the chaincode

To create a k8s chaincode package file, start by creating an `image.json` file.
For example,

```shell
cat << IMAGEJSON-EOF > image.json
{
  "name": "ghcr.io/hyperledger-labs/go-contract",
  "digest": "sha256:802c336235cc1e7347e2da36c73fa2e4b6437cfc6f52872674d1e23f23bba63b"
}
IMAGEJSON-EOF
```

Note: the k8s chaincode package file uses digests because these are immutable, unlike tags.
The docker inspect command can be used to find the digest if required.

```
docker pull ghcr.io/hyperledger-labs/go-contract:v0.7.2
docker inspect --format='{{index .RepoDigests 0}}' ghcr.io/hyperledger-labs/go-contract:v0.7.2 | cut -d'@' -f2
```

Create a `code.tar.gz` archive containing the `image.json` file.

```shell
tar -czf code.tar.gz image.json
```

Create a `metadata.json` file for the chaincode package.
For example,

```shell
cat << METADATAJSON-EOF > metadata.json
{
    "type": "k8s",
    "label": "go-contract"
}
METADATAJSON-EOF
```

Create the final chaincode package archive.

```shell
tar -czf go-contract.tgz metadata.json code.tar.gz
```

Ideally the chaincode package should be created in the same CI/CD pipeline which builds the docker image.
There is an example [package-k8s-chaincode-action](https://github.com/hyperledgendary/package-k8s-chaincode-action) GitHub Action which can create the required k8s chaincode package.

The GitHub Action repository includes a basic shell script which can also be used for automating the process above outside GitHub workflows.
For example, to create a basic k8s chaincode package using the `pkgk8scc.sh` helper script.

```shell
curl -fsSL https://raw.githubusercontent.com/hyperledgendary/package-k8s-chaincode-action/main/pkgk8scc.sh -o pkgk8scc.sh && chmod u+x pkgk8scc.sh
./pkgk8scc.sh -l go-contract -n ghcr.io/hyperledger-labs/go-contract -d sha256:802c336235cc1e7347e2da36c73fa2e4b6437cfc6f52872674d1e23f23bba63b
```

## Chaincode deploy

Deploy the chaincode package as usual, starting by installing the k8s chaincode package.

```shell
peer lifecycle chaincode install go-contract.tgz
```

You can also user the `peer` command to get the chaincode package ID.

```shell
export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid go-contract.tgz) && echo $PACKAGE_ID
```
