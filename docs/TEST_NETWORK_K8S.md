# Kubernetes Test Network

The [Kube Test Network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-k8s) includes support 
for the k8s builder by setting the `TEST_NETWORK_CHAINCODE_BUILDER="k8s"` environment variable.

## Create a Sample Network 

In the `fabric-samples/test-network-k8s` directory:

```shell
export TEST_NETWORK_K8S_CHAINCODE_BUILDER_VERSION="v0.4.0"   # (optional - defaults to v0.4.0)
export TEST_NETWORK_CHAINCODE_BUILDER="k8s"

./network kind 
./network cluster init
```

```shell
./network up
./network channel create
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

The [conga-nft-contract](https://github.com/hyperledgendary/conga-nft-contract) sample chaincode publishes a 
Docker image _and_ a chaincode package archive to GitHub for use with the k8s-builder.  Use of a pre-generated package .tgz 
greatly simplifies the deployment, aligning with standard industry practices for CI/CD and git-ops workflows. 

Download the sample chaincode package: 

```shell
curl -fsSL \
  https://github.com/hyperledgendary/conga-nft-contract/releases/download/v0.1.1/conga-nft-contract-v0.1.1.tgz \
  -o conga-nft-contract-v0.1.1.tgz
```

## Deploying chaincode

Install the chaincode archive to a peer and infer the chaincode's PACKAGE_ID: 

```shell
peer lifecycle chaincode install conga-nft-contract-v0.1.1.tgz
```

```shell
export PACKAGE_ID=$(peer lifecycle chaincode calculatepackageid conga-nft-contract-v0.1.1.tgz) && echo $PACKAGE_ID
```

(Note: PACKAGE_ID must match the chaincode identifier displayed by the `peer lifecycle chaincode install` command.)


Approve the chaincode:

```shell
peer lifecycle \
  chaincode approveformyorg \
  --channelID     mychannel \
  --name          conga-nft-contract \
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
  --name          conga-nft-contract \
  --version       1 \
  --sequence      1 \
  --orderer       org0-orderer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443 \
  --tls --cafile  ${PWD}/build/channel-msp/ordererOrganizations/org0/orderers/org0-orderer1/tls/signcerts/tls-cert.pem
```

## Running transactions

Query the chaincode metadata!

```shell
./network chaincode query conga-nft-contract '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```

## Reset 

Invariably, something in the recipe above will go awry.  Look for additional diagnostics in `network-debug.log` and...

Reset the stage with: 

```shell
./network down && ./network up && ./network channel create
```