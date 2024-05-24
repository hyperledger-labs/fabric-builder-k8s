# Chaincode package

From version 2.0, [Hyperledger Fabric chaincode packages](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html#chaincode-packages) are `.tgz` files which contain two files:

- `metadata.json` — the chaincode label and type
- `code.tar.gz` — source artifacts for the chaincode

Chaincode packages are used by the `peer lifecycle chaincode` command as part of the Fabric chaincode lifecycle to deploy chaincode. For example,

```shell
peer lifecycle chaincode install go-contract.tgz
```

For more information, see the [Fabric chaincode lifecycle](https://hyperledger-fabric.readthedocs.io/en/latest/chaincode_lifecycle.html) documentation.

## metadata.json

The k8s builder will detect chaincode packages which have a type of `k8s`. For example,

```json
{
    "type": "k8s",
    "label": "go-contract"
}
```

The k8s builder uses the chaincode label to label Kubernetes objects, so it must be a [valid Kubernetes label value](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#syntax-and-character-set).

## code.tar.gz

Unlike other chaincode packages, the source artifacts in a k8s chaincode package do not contain the chaincode source files.
Instead, the `code.tar.gz` file contains an `image.json` file which defines which chaincode image should be used.

The `code.tar.gz` file can also contain CouchDB indexes. For more information, see the [CouchDB indexes](https://hyperledger-fabric.readthedocs.io/en/latest/couchdb_as_state_database.html#couchdb-indexes) Fabric documentation.

## image.json

The chaincode image must be built and published before creating the `image.json` file. The `image.json` contains the chaincode image name, and the immutable digest of the published image. For more information, see [Pull an image by digest (immutable identifier)](https://docs.docker.com/engine/reference/commandline/pull/#pull-an-image-by-digest-immutable-identifier). For example.

```json
{
  "name": "ghcr.io/hyperledger-labs/go-contract",
  "digest": "sha256:802c336235cc1e7347e2da36c73fa2e4b6437cfc6f52872674d1e23f23bba63b"
}
```
