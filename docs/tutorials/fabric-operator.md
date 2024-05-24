# Fabric Operator

The [Fabric Operator](https://github.com/hyperledger-labs/fabric-operator) includes support for the k8s builder through its `TEST_NETWORK_PEER_IMAGE` and `TEST_NETWORK_PEER_IMAGE_LABEL` environment variables.

## Create a Sample Network 

Before following [the Fabric Operator sample network instructions](https://github.com/hyperledger-labs/fabric-operator/tree/main/sample-network), export the following environment variables to use the k8s builder peer image:

```shell
export TEST_NETWORK_PEER_IMAGE=ghcr.io/hyperledger-labs/k8s-fabric-peer
export TEST_NETWORK_PEER_IMAGE_LABEL=v0.6.0
```

To create a kind-based sample network using a [fabric-devenv](https://github.com/hyperledgendary/fabric-devenv) VM, run the following commands in the `fabric-operator/sample-network` directory:

```shell
export PATH=$PWD:$PWD/bin:$PATH
export TEST_NETWORK_KUBE_DNS_DOMAIN=test-network
export TEST_NETWORK_INGRESS_DOMAIN=localho.st
network kind
network cluster init
network up
network channel create
```

See the [full Fabric Operator sample network guide](https://github.com/hyperledger-labs/fabric-operator/tree/main/sample-network#k8s-chaincode-builder) for more details, prereqs, and alternative cluster options.

## Set the `peer` CLI environment

Set the `peer` command environment, e.g. for org1, peer1, run the following commands in the `fabric-operator/sample-network` directory: 

```shell
export FABRIC_CFG_PATH=${PWD}/temp/config
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_ADDRESS=test-network-org1-peer1-peer.${TEST_NETWORK_INGRESS_DOMAIN}:443
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_MSPCONFIGPATH=${PWD}/temp/enrollments/org1/users/org1admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/temp/channel-msp/peerOrganizations/org1/msp/tlscacerts/tlsca-signcert.pem
```

## Download a chaincode package

The [sample contracts for Go, Java, and Node.js](https://github.com/hyperledger-labs/fabric-builder-k8s/tree/main/samples) publish a Docker image which the k8s builder can use _and_ a chaincode package file which can be used with the `peer lifecycle chaincode install` command.
Use of a pre-generated chaincode package .tgz greatly simplifies the deployment, aligning with standard industry practices for CI/CD and git-ops workflows.

Download a sample chaincode package, e.g. for the Go contract: 

```shell
curl -fsSL \
  https://github.com/hyperledger-labs/fabric-builder-k8s/releases/download/v0.7.2/go-contract-v0.7.2.tgz \
  -o go-contract-v0.7.2.tgz
```

## Deploying chaincode

Deploy the chaincode package as usual, starting by installing the k8s chaincode package.

```shell
peer lifecycle chaincode install go-contract-v0.7.2.tgz
```

Export a `PACKAGE_ID` environment variable for use in the following commands.

```shell
export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid go-contract-v0.7.2.tgz) && echo $PACKAGE_ID
```

Note: the `PACKAGE_ID` must match the chaincode code package identifier shown by the `peer lifecycle chaincode install` command.

Approve the chaincode:

```shell
peer lifecycle \
  chaincode approveformyorg \
  --channelID     mychannel \
  --name          sample-contract \
  --version       1 \
  --package-id    ${PACKAGE_ID} \
  --sequence      1 \
  --orderer       test-network-org0-orderersnode1-orderer.${TEST_NETWORK_INGRESS_DOMAIN}:443 \
  --tls --cafile  ${PWD}/temp/channel-msp/ordererOrganizations/org0/orderers/org0-orderersnode1/tls/signcerts/tls-cert.pem \
  --connTimeout   15s
```

Commit the chaincode.

```shell
peer lifecycle \
  chaincode commit \
  --channelID     mychannel \
  --name          sample-contract \
  --version       1 \
  --sequence      1 \
  --orderer       test-network-org0-orderersnode1-orderer.${TEST_NETWORK_INGRESS_DOMAIN}:443 \
  --tls --cafile  ${PWD}/temp/channel-msp/ordererOrganizations/org0/orderers/org0-orderersnode1/tls/signcerts/tls-cert.pem \
  --connTimeout   15s
```

Inspect chaincode pods.

```shell
kubectl -n test-network describe pods -l app.kubernetes.io/created-by=fabric-builder-k8s
```

## Running transactions

Query the chaincode metadata!

```shell
network chaincode query sample-contract '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```
