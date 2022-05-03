# fabric-builder-k8s

Proof of concept Fabric builder for Kubernetes

Status: it should just about work now but there are a few issues to iron out (and tests to write) before it's properly usable!

## Usage

**Note:** See [Kubernetes Test Network](docs/TEST_NETWORK_K8S.md) for specific instructions for using the builder with the k8s test network.

The k8s builder can be run in cluster using the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` environment variables, or it can connect using a `KUBECONFIG_PATH` environment variable.

An optional `KUBE_NAMESPACE` can be used to specify the namespace to deploy chaincode to.

A `CORE_PEER_ID` environment variable is also currently required.

External builders are configured in the `core.yaml` file, for example:

```
  externalBuilders:
    - name: k8s_builder
      path: /home/peer/ccbuilders/k8s_builder
      propagateEnvironment:
        - CORE_PEER_ID
        - KUBE_NAMESPACE
        - KUBERNETES_SERVICE_HOST
        - KUBERNETES_SERVICE_PORT
```

See [External Builders and Launchers](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) for details of Hyperledger Fabric builders.

## Chaincode package

The k8s chaincode package must contain an image name and tag.
For example, to create a basic k8s chaincode package using the `pkgk8scc.sh` helper script.

```shell
pkgk8scc.sh -l conga-nft-contract -n ghcr.io/hyperledgendary/conga-nft-contract -t b96d4701d6a04e6109bc51ef1c148a149bfc6200
```

You can also create the same chaincode package manually.
Start by creating an `image.json` file.

```shell
cat << IMAGEJSON-EOF > image.json
{
  "name": "ghcr.io/hyperledgendary/conga-nft-contract",
  "tag": "b96d4701d6a04e6109bc51ef1c148a149bfc6200"
}
IMAGEJSON-EOF
```

Create a `code.tar.gz` archive containing the `image.json` file.

```shell
tar -czf code.tar.gz image.json
```

Create a `metadata.json` file for the chaincode package.

```shell
cat << METADATAJSON-EOF > metadata.json
{
    "type": "k8s",
    "label": "conga-nft-contract"
}
METADATAJSON-EOF
```

Create the final chaincode package archive.

```shell
tar -czf conga-nft-contract.tgz metadata.json code.tar.gz
```
## Chaincode deploy

Deploy the chaincode package as usual, starting by installing the k8s chaincode package.

```shell
peer lifecycle chaincode install conga-nft-contract.tgz
```
