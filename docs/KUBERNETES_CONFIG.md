# Kubernetes Configuration

## Builder requirements

The k8s builder needs sufficient permissions to manage chaincode pods on behalf of the Fabric `peer`.

| Resource | Permissions                      |
| -------- | -------------------------------- |
| pods     | get, list, watch, create, delete |
| secrets  | create, patch                    |

For example, follow these steps if the builder will be running in the `default` namespace using the `default` service account.

1. Create a `fabric-builder-role` role.

```shell
cat <<EOF | kubectl apply -f -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: fabric-builder-role
  namespace: default
rules:
  - apiGroups:
      - ""
      - apps
    resources:
      - pods
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
```

2. Create a `fabric-builder-rolebinding` role binding.

```shell
cat <<EOF | kubectl apply -f -
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: fabric-builder-rolebinding
  namespace: default
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: fabric-builder-role 
subjects:
- namespace: default 
  kind: ServiceAccount
  name: default 
EOF
```

3. Check the permissions, e.g.

```shell
kubectl auth can-i patch secrets --namespace default --as system:serviceaccount:default:default
```

## Chaincode requirements

By default, the k8s builder starts chaincode pods in the same namespace, however the `FABRIC_K8S_BUILDER_NAMESPACE` environment variable can be used to start chaincode pods in a different namespace.

Chaincode pods are created with a service account defined by the `FABRIC_K8S_BUILDER_SERVICE_ACCOUNT` environment variable, or the `default` service account if the variable is not set.
If your chaincode images are published to registries which require credentials, you will need to add image pull secrets to the service account.

For example, follow these steps if `FABRIC_K8S_BUILDER_NAMESPACE` and `FABRIC_K8S_BUILDER_SERVICE_ACCOUNT` are both set to `fabric-chaincode`.

1. Create the `fabric-chaincode` namespace.

```shell
kubectl create namespace chaincode
```

2. Create the `fabric-chaincode` service account.

```shell
kubectl create serviceaccount fabric-chaincode --namespace=fabric-chaincode
```

3. Create an imagePullSecret.

```shell
kubectl create secret docker-registry fabregistrykey --namespace=fabric-chaincode \
    --docker-server=DOCKER_SERVER \
    --docker-username=DOCKER_USERNAME \
    --docker-password=DOCKER_PASSWORD \
    --docker-email=DOCKER_EMAIL
```

4. Add the image pull secret to the service account.

```shell
kubectl patch serviceaccount fabric-chaincode --namespace=fabric-chaincode -p '{"imagePullSecrets": [{"name": "fabregistrykey"}]}'
```

See the Kubernetes [Configure Service Accounts for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#add-imagepullsecrets-to-a-service-account) documentation for details.
