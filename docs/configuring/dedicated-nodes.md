# Dedicated nodes

By default, the k8s builder does not implement any Kubernetes node scheduling strategies.

The `FABRIC_K8S_BUILDER_NODE_ROLE` environment variable can be used to schedule chaincode on dedicated Kubernetes nodes.
Chaincode pods will be configured with an affinity for nodes with the `fabric-builder-k8s-role=<node_role>` label, and will tolerate nodes with the `fabric-builder-k8s-role=<node_role>:NoSchedule` taint.

For example, if `FABRIC_K8S_BUILDER_NODE_ROLE` is set to `chaincode`, use the following `kubectl` commands to configure a dedicated chaincode node `ccnode`.

```shell
kubectl label nodes ccnode fabric-builder-k8s-role=chaincode
kubectl taint nodes ccnode fabric-builder-k8s-role=chaincode:NoSchedule
```

More complex requirements should be handled with Dynamic Admission Control using a Mutating Webhook.
For example, you could use a webhook to assign node affinity and tolerations to all pods in a `chaincode` namespace.
