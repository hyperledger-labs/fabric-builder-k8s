# How to contribute

There are lots of ways to contribute, and we :heart: them all!

## Questions

If you have any suggestions, questions about the project, or are interested in contributing, you can find us in the [fabric-kubernetes](https://discord.com/channels/905194001349627914/945796983795384331) channel on Discord.

## Issues

All issues are tracked in [our issues tab in GitHub](https://github.com/hyperledger-labs/fabric-builder-k8s/issues). If you find a bug which we don't already know about, you can help us by creating a new issue describing the problem. Please include as much detail as possible to help us track down the cause.

## Fixes

If you want to begin contributing code, looking through our open issues is a good way to start. Try looking for recent issues with detailed descriptions first.

## Enhancements

Open an issue to make sure your contibution is likely to be accepted before investing a lot of effort in larger changes.

## Pull Requests

We use [pull requests](http://help.github.com/pull-requests/) to deliver changes to the code. Follow these steps to deliver your first pull request:

1. [Fork the repository](https://guides.github.com/activities/forking/#fork) and create a new branch from `main`.
2. If you've added code that should be tested, add tests!
3. If you've added any new features or made breaking changes, update the documentation.
4. Ensure all the tests pass.
5. Include a descriptive message, and the [Developer Certificate of Origin (DCO) sign-off](https://github.com/probot/dco#how-it-works) on all commit messages.
6. [Issue a pull request](https://guides.github.com/activities/forking/#making-a-pull-request)!
7. [GitHub Actions](https://github.com/hyperledger-labs/fabric-builder-k8s/actions) builds must succeed before the pull request can be reviewed and merged.

## Coding Style

Please to try to be consistent with the rest of the code and conform to linting rules where they are provided.

## Development environment

There is a [Visual Studio Code Dev Container](https://code.visualstudio.com/docs/devcontainers/containers) which should help develop and test the k8s builder in a consistent development environment.
It includes a preconfigured nano Fabric test network and minikube which can be used to run end to end tests.

Build your latest k8s builder changes.

```
GOBIN="${PWD}"/.fabric/builders/k8s_builder/bin go install ./cmd/...
```

[Configure kubernetes](./docs/KUBERNETES_CONFIG.md) and export the kubeconfig path.

```
export KUBECONFIG_PATH="${HOME}/.kube/config"
```

Start the Fabric test network in the `.fabric/test-network-nano-bash` directory.

```
./network.sh start
```

In a new shell in the `.fabric/test-network-nano-bash` directory.

```shell
curl -fsSL \
  https://github.com/hyperledger-labs/fabric-builder-k8s/releases/download/v0.7.2/go-contract-v0.7.2.tgz \
  -o go-contract-v0.7.2.tgz
```

Set up the environment for running peer commands and check everything is working.

```
. ./peer1admin.sh
peer channel list
```

Deploy the chaincode package as usual, starting by installing the k8s chaincode package.

```shell
peer lifecycle chaincode install go-contract-v0.7.2.tgz
```

Export a `PACKAGE_ID` environment variable for use in the following commands.

```shell
export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid go-contract-v0.7.2.tgz) && echo $PACKAGE_ID
```

Note: the `PACKAGE_ID` must match the chaincode code package identifier shown by the `peer lifecycle chaincode install` command.

Approve the chaincode.

```shell
peer lifecycle chaincode approveformyorg -o 127.0.0.1:6050 --channelID mychannel --name sample-contract --version 1 --package-id $PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

Commit the chaincode.

```shell
peer lifecycle chaincode commit -o 127.0.0.1:6050 --channelID mychannel --name sample-contract --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

Query the chaincode metadata!

```shell
peer chaincode query -C mychannel -n sample-contract -c '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```

Debug chaincode!!

## Code of Conduct Guidelines <a name="conduct"></a>

See our [Code of Conduct Guidelines](./CODE_OF_CONDUCT.md).

