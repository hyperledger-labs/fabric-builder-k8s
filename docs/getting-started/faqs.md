# Frequently Asked Questions

## Chaincode

### Are private chaincode images supported?

Yes. For more information, see [Kubernetes service account](../configuring/kubernetes-service-account.md)

### Will chaincode work in multi-architecture Fabric networks?

Yes. Build multi-architecture chaincode images if you have a multi-architecture Fabric network.

### Can every chaincode be launched in a different namespace?

No. The k8s builder configuration is defined for each peer, so all chaincode launched by the same peer will run in the same namespace.
