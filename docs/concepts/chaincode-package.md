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
