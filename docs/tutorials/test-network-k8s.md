# Kubernetes Test Network

The [Kube Test Network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-k8s) includes support for the k8s builder by setting the `TEST_NETWORK_CHAINCODE_BUILDER="k8s"` environment variable.

## Create a Sample Network 

In the `fabric-samples/test-network-k8s` directory:

```shell
export PATH=$PWD:$PWD/bin:$PATH

export TEST_NETWORK_K8S_CHAINCODE_BUILDER_VERSION="v0.6.0"   # (optional - defaults to v0.4.0)
export TEST_NETWORK_CHAINCODE_BUILDER="k8s"

network kind 
network cluster init
network up
network channel create
```

(Check / follow the detailed log file for errors and progress at `network-debug.log`.  E.g. in a separate shell:)
```shell
tail -f network-debug.log
```

## Set the `peer` CLI environment

Make sure the `build` directory exists -- this will be created by `network channel create`. 

Set the `peer` command environment, e.g. for org1, peer1: 

```shell
export FABRIC_CFG_PATH=${PWD}/config/org1
export CORE_PEER_ADDRESS=org1-peer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443
export CORE_PEER_MSPCONFIGPATH=${PWD}/build/enrollments/org1/users/org1admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/build/channel-msp/peerOrganizations/org1/msp/tlscacerts/tlsca-signcert.pem
```

## Download a chaincode package

The [sample contracts for Go, Java, and Node.js](samples/README.md) publish a Docker image which the k8s builder can use _and_ a chaincode package file which can be used with the `peer lifecycle chaincode install` command.
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
  --orderer       org0-orderer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443 \
  --tls --cafile  ${PWD}/build/channel-msp/ordererOrganizations/org0/orderers/org0-orderer1/tls/signcerts/tls-cert.pem
```

Commit the chaincode.

```shell
peer lifecycle \
  chaincode commit \
  --channelID     mychannel \
  --name          sample-contract \
  --version       1 \
  --sequence      1 \
  --orderer       org0-orderer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443 \
  --tls --cafile  ${PWD}/build/channel-msp/ordererOrganizations/org0/orderers/org0-orderer1/tls/signcerts/tls-cert.pem
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

## Reset 

Invariably, something in the recipe above will go awry.  Look for additional diagnostics in `network-debug.log` and...

Reset the stage with: 

```shell
network down && network up && network channel create
```
