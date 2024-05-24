# Kubernetes permissions

The k8s builder needs sufficient permissions to manage chaincode pods on behalf of the Fabric `peer`.

| Resource | Permissions                      |
| -------- | -------------------------------- |
| jobs     | get, list, watch, create         |
| pods     | get, list, watch, create, delete |
| secrets  | create, patch                    |

For example, follow these steps if the builder will be running in the `default` namespace using the `default` service account.

Create a `fabric-builder-role` role.

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
      - batch
    resources:
      - jobs
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

Create a `fabric-builder-rolebinding` role binding.

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

Check the permissions, e.g.

```shell
kubectl auth can-i patch secrets --namespace default --as system:serviceaccount:default:default
```
