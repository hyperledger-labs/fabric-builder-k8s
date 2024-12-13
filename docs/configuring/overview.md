# Configuration overview

Fabric peers must be configured to use the k8s external builder, and to propagate the required environment variables to configure the builder.

## Fabric peer configuration

External builders are configured in the `core.yaml` file, for example:

```
  externalBuilders:
    - name: k8s_builder
      path: /opt/hyperledger/k8s_builder
      propagateEnvironment:
        - CORE_PEER_ID
        - FABRIC_K8S_BUILDER_DEBUG
        - FABRIC_K8S_BUILDER_NAMESPACE
        - FABRIC_K8S_BUILDER_NODE_ROLE
        - FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX
        - FABRIC_K8S_BUILDER_SERVICE_ACCOUNT
        - FABRIC_K8S_BUILDER_START_TIMEOUT
        - KUBERNETES_SERVICE_HOST
        - KUBERNETES_SERVICE_PORT
```

For more information, see [Configuring external builders and launchers](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html#configuring-external-builders-and-launchers) in the Fabric documentation.

## Environment variables

The k8s builder is configured using the following environment variables.

| Name                                  | Default                          | Description                                          |
| ------------------------------------- | -------------------------------- | ---------------------------------------------------- |
| CORE_PEER_ID                          |                                  | The Fabric peer ID (required)                        |
| FABRIC_K8S_BUILDER_NAMESPACE          | The peer namespace or `default`  | The Kubernetes namespace to run chaincode with       |
| FABRIC_K8S_BUILDER_NODE_ROLE          |                                  | Use dedicated Kubernetes nodes to run chaincode      |
| FABRIC_K8S_BUILDER_OBJECT_NAME_PREFIX | `hlfcc`                          | Eye-catcher prefix for Kubernetes object names       |
| FABRIC_K8S_BUILDER_SERVICE_ACCOUNT    | `default`                        | The Kubernetes service account to run chaincode with |
| FABRIC_K8S_BUILDER_START_TIMEOUT      | `3m`                             | The timeout when waiting for chaincode pods to start |
| FABRIC_K8S_BUILDER_DEBUG              | `false`                          | Set to `true` to enable k8s builder debug messages   |

The k8s builder can be run in cluster using the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` environment variables, or it can connect using a `KUBECONFIG_PATH` environment variable.
