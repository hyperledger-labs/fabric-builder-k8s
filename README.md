# fabric-builder-k8s

Proof of concept Fabric builder for Kubernetes


## Deployment

Set up a [k8s test network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-k8s):

1. Configure the external builder in kube/org[1,2]/core.yaml:
```
    externalBuilders:
      - name: k8s_builder
        path: /var/hyperledger/fabric/external_builders/k8s_builder
        propagateEnvironment:
          - CHAINCODE_AS_A_SERVICE_BUILDER_CONFIG
```
2. `network kind`
3. `network up && network channel create`

Deploy the fabric-builder-k8s.  

This will create a role-binding, allowing the default service account to create 
Pods, Deployments, and Services in the kube namespace.  In addition, the kustomization creates two jobs to 
copy the k8s_builder binaries from the Docker image into the org-specific persistent volume mount.
```
kubectl kustomize config/default | kubectl apply -f - 
```

## Docker image 

```
docker build -t ghcr.io/hyperledgendary/fabric-builder-k8s . 

kind load docker-image ghcr.io/hyperledgendary/fabric-builder-k8s 
```

## Chaincode package

The k8s chaincode package contains an image name and tag.
For example, to create a basic k8s chaincode package using the `pkgk8scc.sh` helper script.

```shell
pkgk8scc.sh -l sample -n ghcr.io/hyperledgendary/fabric-ccaas-asset-transfer-basic -t latest
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

## Chaincode install

In the `fabric-samples/test-network-k8s` directory...

Configure the `peer` command environment, e.g. for org1, peer1

```shell
export FABRIC_CFG_PATH=${PWD}/config/org1
export CORE_PEER_ADDRESS=org1-peer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443
export CORE_PEER_MSPCONFIGPATH=${PWD}/build/enrollments/org1/users/org1admin/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/build/channel-msp/peerOrganizations/org1/msp/tlscacerts/tlsca-signcert.pem
```

Install the chaincode package

```shell
peer lifecycle chaincode install sample.tgz
```

```shell
export PACKAGE_ID=sample:...
```

Approve the chaincode

```shell
peer lifecycle \
  chaincode approveformyorg \
  --channelID     mychannel \
  --name          sample \
  --version       1 \
  --package-id    ${PACKAGE_ID} \
  --sequence      1 \
  --orderer       org0-orderer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443 \
  --tls --cafile  ${PWD}/build/channel-msp/ordererOrganizations/org0/orderers/org0-orderer1/tls/signcerts/tls-cert.pem
```

Commit the chaincode

```
peer lifecycle \
  chaincode commit \
  --channelID     mychannel \
  --name          sample \
  --version       1 \
  --sequence      1 \
  --orderer       org0-orderer1.${TEST_NETWORK_DOMAIN:-vcap.me}:443 \
  --tls --cafile  ${PWD}/build/channel-msp/ordererOrganizations/org0/orderers/org0-orderer1/tls/signcerts/tls-cert.pem
```

Configure the `network` script to use the new chaincode

```shell
export TEST_NETWORK_CHAINCODE_NAME=sample
```

Query the chaincode metadata!

```shell
./network chaincode query '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```
