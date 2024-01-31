# Hyperledger Fabric Operator

The k8s builder can be used with the [Hyperledger Fabric Operator](https://github.com/hyperledger-labs/hlf-operator) by following the instructions below.

## Create a Demo Network 

Follow the [Hyperledger Fabric Operator getting started instructions](https://labs.hyperledger.org/hlf-operator/docs/getting-started) with the following modifications (TBC)...

Configure the hfl-operator to use the k8s builder peer image:

```shell
export PEER_IMAGE=ghcr.io/hyperledger-labs/k8s-fabric-peer
export PEER_VERSION=v0.6.0
```

After creating the peer, patch it to include the k8s builder configuration:

```shell
kubectl patch peer org1-peer0 --type=json --patch-file=/dev/stdin <<-EOF
[
  {
    "op" : "add",
    "path" : "/spec/externalBuilders/-",
    "value" : {
      "name" : "k8s_builder",
      "path" : "/opt/hyperledger/k8s_builder",
      "propagateEnvironment" : [
        "CORE_PEER_ID",
        "FABRIC_K8S_BUILDER_DEBUG",
        "FABRIC_K8S_BUILDER_NAMESPACE",
        "FABRIC_K8S_BUILDER_SERVICE_ACCOUNT",
        "KUBERNETES_SERVICE_HOST",
        "KUBERNETES_SERVICE_PORT"
      ]
    }
  }
]
EOF
```

Note: the configuration change does not get picked up without restarting the pods:

```shell
kubectl scale deployment org1-peer0 --replicas=0 -n default
kubectl scale deployment org1-peer0 --replicas=1 -n default
```

TODO: Is there a better way to create peers with the builder pre-configured?

Ensure the k8s builder has the required permissions to manage chaincode pods:

```shell
cat <<EOF | kubectl apply -f -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: k8s-builder-role
  namespace: default
rules:
  - apiGroups:
      - ""
      - apps
    resources:
      - pods
      - deployments
      - configmaps
      - secrets
    verbs:
      - get
      - list
      - watch
      - create
      - delete
      - patch
EOF

cat <<EOF | kubectl apply -f -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: k8s-builder-rolebinding
  namespace: default 
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: k8s-builder-role 
subjects:
- namespace: default 
  kind: ServiceAccount
  name: org1-peer0 
EOF

kubectl auth can-i list pods --namespace default --as system:serviceaccount:default:org1-peer0
kubectl auth can-i delete pods --namespace default --as system:serviceaccount:default:org1-peer0
kubectl auth can-i patch secrets --namespace default --as system:serviceaccount:default:org1-peer0
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

Install the chaincode package to a peer

**Note:** the `--language` argument is required even though it does not make sense in this case

```shell
export CHAINCODE_NAME=go-contract
export CHAINCODE_LABEL=go-contract

kubectl hlf chaincode install --path=./go-contract-v0.7.2.tgz \
    --config=org1.yaml --language=golang --label=$CHAINCODE_LABEL --user=admin --peer=org1-peer0.default
```

Get the chaincode's PACKAGE_ID

```
export PACKAGE_ID=$(kubectl hlf chaincode calculatepackageid --path=./go-contract-v0.7.2.tgz --language=golang --label=$CHAINCODE_LABEL) && echo $PACKAGE_ID
```

Approve the chaincode

```shell
export SEQUENCE=1
export VERSION="1.0"
kubectl hlf chaincode approveformyorg --config=org1.yaml --user=admin --peer=org1-peer0.default \
    --package-id=$PACKAGE_ID \
    --version "$VERSION" --sequence "$SEQUENCE" --name=$CHAINCODE_NAME \
    --policy="OR('Org1MSP.member')" --channel=demo
```

Commit the chaincode

```shell
kubectl hlf chaincode commit --config=org1.yaml --user=admin --mspid=Org1MSP \
    --version "$VERSION" --sequence "$SEQUENCE" --name=$CHAINCODE_NAME \
    --policy="OR('Org1MSP.member')" --channel=demo
```

## Running transactions

Query the chaincode metadata!

```shell
kubectl hlf chaincode query --config=org1.yaml \
    --user=admin --peer=org1-peer0.default \
    --chaincode=$CHAINCODE_NAME --channel=demo \
    --fcn=org.hyperledger.fabric:GetMetadata -a '[]'
```
