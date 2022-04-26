# fabric-builder-k8s

Proof of concept Fabric builder for Kubernetes

## Deploying with k8s test network

The k8s builder can be used with the [k8s test network](https://github.com/hyperledger/fabric-samples/tree/main/test-network-k8s).

Before following the instructions to set up the k8s test network, it needs to be configured to use the k8s builder peer image.
Find the latest [k8s-fabric-peer](https://github.com/hyperledgendary/fabric-builder-k8s/pkgs/container/k8s-fabric-peer) image and export a `TEST_NETWORK_FABRIC_PEER_IMAGE` environment variable, e.g.

```shell
export TEST_NETWORK_FABRIC_PEER_IMAGE=ghcr.io/hyperledgendary/k8s-fabric-peer:47ec271bb9d7b31f35bcb5f0bd499835a223c5c6
```

The org1 and org2 `core.yaml` files also need to be updated with the k8s builder configuration.

```
  externalBuilders:
    - name: k8s_builder
      path: /opt/hyperledger/k8s_builder
      propagateEnvironment:
        - KUBERNETES_SERVICE_HOST
        - KUBERNETES_SERVICE_PORT
```

You can use [yq](https://mikefarah.gitbook.io/yq/) to update the `core.yaml` files.
Make sure you are in the `fabric-samples/test-network-k8s` directory before running the following commands.

```shell
yq -i '.chaincode.externalBuilders += { "name": "k8s_builder", "path": "/opt/hyperledger/k8s_builder", "propagateEnvironment": [ "KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT" ] }' config/org1/core.yaml
yq -i '.chaincode.externalBuilders += { "name": "k8s_builder", "path": "/opt/hyperledger/k8s_builder", "propagateEnvironment": [ "KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT" ] }' config/org2/core.yaml
```

After launching the k8s test network using the `./network up` command, you also need to configure a k8s service user role to allow the k8s builder to create chaincode deployments.

TODO: Create a role (cut this down to what is actually required!)

```shell
cat <<EOF | kubectl apply -f -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: fabric-builder-role
  namespace: test-network
rules:
  - apiGroups:
        - ""
        - apps
        - autoscaling
        - batch
        - extensions
        - policy
        - rbac.authorization.k8s.io
    resources:
      - pods
      - componentstatuses
      - configmaps
      - daemonsets
      - deployments
      - events
      - endpoints
      - horizontalpodautoscalers
      - ingress
      - jobs
      - limitranges
      - namespaces
      - nodes
      - pods
      - persistentvolumes
      - persistentvolumeclaims
      - resourcequotas
      - replicasets
      - replicationcontrollers
      - serviceaccounts
      - services
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
EOF
```

Create a role binding.

```shell
cat <<EOF | kubectl apply -f -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: fabric-builder-rolebinding
  namespace: test-network 
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: fabric-builder-role 
subjects:
- namespace: test-network 
  kind: ServiceAccount
  name: default 
EOF
```

And finally, check it worked!

```shell
kubectl auth can-i create deployments --namespace test-network --as system:serviceaccount:test-network:default
```

## Chaincode package

The k8s chaincode package contains an image name and tag.
For example, to create a basic k8s chaincode package using the `pkgk8scc.sh` helper script.

```shell
pkgk8scc.sh -l sample-k8s-cc -n ghcr.io/hyperledger/asset-transfer-basic -t 1.0
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
tar -czf sample-k8s-cc.tgz metadata.json code.tar.gz
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
peer lifecycle chaincode install sample-k8s-cc.tgz
```

Export a `PACKAGE_ID` environment variable for use in the following commands

```shell
export PACKAGE_ID=sample:$(shasum -a 256 sample-k8s-cc.tgz  | tr -s ' ' | cut -d ' ' -f 1)
```

Note: this should match the chaincode code package identifier shown by the `peer lifecycle chaincode install` command

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

```shell
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
