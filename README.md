# fabric-builder-k8s

Proof of concept Fabric builder for Kubernetes

## Chaincode package

The k8s chaincode package contains an image name and tag.
For example, to create a basic k8s chaincode package using the `pkgk8scc.sh` helper script.

```shell
pkgk8scc.sh -l sample -n ghcr.io/hyperledger/asset-transfer-basic -t 1.0
```

You can also create the same chaincode package manually.
Start by creating an `image.json` file.

```shell
cat << IMAGEJSON-EOF > image.json
{
  "name": "ghcr.io/hyperledger/asset-transfer-basic",
  "tag": "1.0"
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
    "label": "sample"
}
METADATAJSON-EOF
```

Create the final chaincode package archive.

```shell
tar -czf sample-k8s-cc.tar.gz metadata.json code.tar.gz
```
