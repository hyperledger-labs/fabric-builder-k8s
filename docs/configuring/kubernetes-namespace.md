# Kubernetes namespace

By default, the k8s builder starts chaincode pods in the same namespace as the peer, or the `default` namespace if the peer is running outside Kubernetes.

The `FABRIC_K8S_BUILDER_NAMESPACE` environment variable can be used to start chaincode pods in a different namespace.

For example, if `FABRIC_K8S_BUILDER_NAMESPACE` is set to `hlf-chaincode`, create the required namespace using the following command.

```shell
kubectl create namespace hlf-chaincode
```
