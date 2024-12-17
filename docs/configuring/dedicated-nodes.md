# Dedicated nodes

TBC

The `FABRIC_K8S_BUILDER_NODE_ROLE` environment variable can be used to...

For example, if `FABRIC_K8S_BUILDER_NODE_ROLE` is set to `chaincode`, ... using the following command.

```shell
kubectl label nodes node1 fabric-builder-k8s-role=chaincode
kubectl taint nodes node1 fabric-builder-k8s-role=chaincode:NoSchedule
```

More complex requirements should be handled with Dynamic Admission Control using a Mutating Webhook.
For example, it looks like the namespace-node-affinity webhook could be used to assign node affinity and tolerations to all pods in the FABRIC_K8S_BUILDER_NAMESPACE namespace.
