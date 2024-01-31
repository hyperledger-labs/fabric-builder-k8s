## Usage

The k8s builder can be run in cluster using the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` environment variables, or it can connect using a `KUBECONFIG_PATH` environment variable.

The following optional environment variables can be used to configure the k8s builder:

- `FABRIC_K8S_BUILDER_DEBUG` whether to enable additional logging
- `FABRIC_K8S_BUILDER_NAMESPACE` specifies the namespace to deploy chaincode to
- `FABRIC_K8S_BUILDER_SERVICE_ACCOUNT` specifies the service account for the chaincode pod

A `CORE_PEER_ID` environment variable is also currently required.

External builders are configured in the `core.yaml` file, for example:

```
  externalBuilders:
    - name: k8s_builder
      path: /opt/hyperledger/k8s_builder
      propagateEnvironment:
        - CORE_PEER_ID
        - FABRIC_K8S_BUILDER_DEBUG
        - FABRIC_K8S_BUILDER_NAMESPACE
        - FABRIC_K8S_BUILDER_SERVICE_ACCOUNT
        - KUBERNETES_SERVICE_HOST
        - KUBERNETES_SERVICE_PORT
```

See [External Builders and Launchers](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) for details of Hyperledger Fabric builders.

As well as configuring Fabric to use the k8s builder, you will need to [configure Kubernetes](docs/KUBERNETES_CONFIG.md) to allow the builder to start chaincode pods successfully.