# Quick start

The [fabric-samples](https://github.com/hyperledger/fabric-samples/) Kubernetes test network includes support for the k8s builder and provides the quickest way to get started.

First create a directory to download all the required files and run the demo.

```shell
mkdir k8s-builder-demo
cd k8s-builder-demo
```

Now follow the steps below to deploy your first smart contract using the k8s builder!

## Download the Kubernetes test network

Download the sample Kubernetes test network (fabric-samples isn't tagged so we'll use a known good commit).

```shell
export FABRIC_SAMPLES_COMMIT=1058f9ffe16add583d1a11342deb5a9df3e5b72c
curl -sSL "https://github.com/hyperledger/fabric-samples/archive/${FABRIC_SAMPLES_COMMIT}.tar.gz" | \
  tar -xzf - --strip-components=1 \
    fabric-samples-${FABRIC_SAMPLES_COMMIT}/test-network-k8s
```

## Configure the Kubernetes test network

Set the following environment variables to enable the k8s builder and define which version to use.

```shell
export TEST_NETWORK_CHAINCODE_BUILDER="k8s"
export TEST_NETWORK_K8S_CHAINCODE_BUILDER_VERSION="0.15.1"
```

## Download chaincode samples

The Kubernetes test network instructions deploy the `asset-transfer-basic` sample. The `asset-transfer-basic` sample may work with the k8s builder in some environments however it is better to download and use the [samples provided with the k8s builder](https://github.com/hyperledger-labs/fabric-builder-k8s/tree/main/samples) instead.

```shell
curl -sSL "https://github.com/hyperledger-labs/fabric-builder-k8s/archive/refs/tags/v${TEST_NETWORK_K8S_CHAINCODE_BUILDER_VERSION}.tar.gz" | \
  tar -xzf - --strip-components=2 fabric-builder-k8s-${TEST_NETWORK_K8S_CHAINCODE_BUILDER_VERSION}/samples
```

## Start the Kubernetes test network

In the `test-network-k8s` directory, follow the [Kubernetes test network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-k8s) instructions to launch the network, and create a channel. Stop before deploying the `asset-transfer-basic` smart contract.

## Deploy a sample contract

Use the Kubernetes test network script to deploy one of the k8s builder's sample contracts.

```shell
./network chaincode deploy sample-contract ../go-contract
```

You can query the chaincode metadata to confirm that the sample was deployed successfully.

```shell
./network chaincode query sample-contract '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```

Use the `kubectl` command to inspect chaincode jobs.

```shell
kubectl -n test-network describe jobs -l app.kubernetes.io/created-by=fabric-builder-k8s
```

## Running transactions

Use the following commands to invoke and query transactions on the sample contract.

```shell
./network chaincode invoke sample-contract '{"Args":["PutValue","asset1","green"]}'
./network chaincode query sample-contract '{"Args":["GetValue","asset1"]}'
```

## Cleaning up

Follow the [Kubernetes test network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-k8s) instructions to tear down the network.
