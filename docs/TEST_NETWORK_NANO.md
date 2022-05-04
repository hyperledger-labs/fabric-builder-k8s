# Nano Test Network

The k8s builder can be used with the [nano test network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-nano-bash) by following the instructions below.

## Download builder binaries

Download the latest builder binaries from the [releases page](https://github.com/hyperledgendary/fabric-builder-k8s/releases) and extract them to a `k8s_builder/bin` directory in your home directory.

## Configure builder

After installing the nano test network prereqs, the `fabric-samples/config/core.yaml` file needs to be updated with the k8s builder configuration.

```
  externalBuilders:
    - name: k8s_builder
      path: <HOME>/k8s_builder
      propagateEnvironment:
        - CORE_PEER_ID
        - KUBECONFIG_PATH
```

You can use [yq](https://mikefarah.gitbook.io/yq/) to update the `fabric-samples/config/core.yaml` files.
Make sure you are in the `fabric-samples` directory before running the following commands.

```shell
FABRIC_K8S_BUILDER_PATH=${HOME}/k8s_builder yq -i '.chaincode.externalBuilders += { "name": "k8s_builder", "path": "${FABRIC_K8S_BUILDER_PATH}" | envsubst(ne), "propagateEnvironment": [ "CORE_PEER_ID", "KUBECONFIG_PATH" ] }' config/core.yaml
```

## Kubernetes configuration

The k8s builder needs a kubeconfig file to access a Kubernetes cluster to deploy chaincode. Make sure the `KUBECONFIG_PATH` environment variable is available on every peer the builder is configured on.

```shell
export KUBECONFIG_PATH=$HOME/.kube/config
```

## Deploying chaincode

Deploy the chaincode package as usual, starting by installing the k8s chaincode package.

```shell
peer lifecycle chaincode install conga-nft-contract.tgz
```

Export a `PACKAGE_ID` environment variable for use in the following commands.

```shell
export PACKAGE_ID=conga-nft-contract:$(shasum -a 256 conga-nft-contract.tgz  | tr -s ' ' | cut -d ' ' -f 1)
```

Note: this should match the chaincode code package identifier shown by the `peer lifecycle chaincode install` command.

Approve the chaincode.

```shell
peer lifecycle chaincode approveformyorg -o 127.0.0.1:6050 --channelID mychannel --name conga-nft-contract --version 1 --package-id $PACKAGE_ID --sequence 1 --tls --cafile ${PWD}/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

Commit the chaincode.

```shell
peer lifecycle chaincode commit -o 127.0.0.1:6050 --channelID mychannel --name conga-nft-contract --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

## Running transactions

Query the chaincode metadata!

```shell
peer chaincode query -C mychannel -n conga-nft-contract -c '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```
