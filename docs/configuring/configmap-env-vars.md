# ConfigMap Environment Variables

The Fabric K8s Builder supports injecting environment variables into chaincode pods from Kubernetes ConfigMaps. This allows you to configure chaincode behavior without rebuilding container images.

## Overview

Environment variables are automatically loaded from a ConfigMap if it exists with the same name as the chaincode label. No additional configuration is required on the peer.

## How It Works

When deploying chaincode:
1. The builder extracts the chaincode label from the package ID
2. It checks if a ConfigMap exists with that exact name
3. If found, all key-value pairs from the ConfigMap are mounted as environment variables
4. If not found, the chaincode pod starts normally without the additional environment variables

## Configuration

### 1. Package Your Chaincode

Package your chaincode with a label:

```bash
peer lifecycle chaincode package mycc.tar.gz \
  --path ./chaincode \
  --lang k8s \
  --label mycc_v1
```

The label `mycc_v1` will be used as the ConfigMap name.

### 2. Create the ConfigMap

Create a ConfigMap with the **exact same name** as your chaincode label:

```bash
kubectl create configmap mycc_v1 \
  --from-literal=LOG_LEVEL=debug \
  --from-literal=DATABASE_URL=postgres://db:5432/mydb \
  --from-literal=API_KEY=secret123 \
  -n hyperledger
```

Or using a YAML file:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: mycc_v1
  namespace: hyperledger
data:
  LOG_LEVEL: "debug"
  DATABASE_URL: "postgres://db:5432/mydb"
  API_KEY: "secret123"
  FEATURE_FLAG_X: "enabled"
```

Apply it:

```bash
kubectl apply -f mycc-configmap.yaml
```

### 3. Deploy Your Chaincode

Install and deploy your chaincode normally:

```bash
peer lifecycle chaincode install mycc.tar.gz
peer lifecycle chaincode approveformyorg --channelID mychannel --name mycc --version 1.0 --package-id mycc_v1:hash...
peer lifecycle chaincode commit --channelID mychannel --name mycc --version 1.0
```

The chaincode pod will automatically have access to all environment variables from the ConfigMap.

## Naming Convention

**Important:** The ConfigMap name must **exactly match** the chaincode label.

### Examples

| Chaincode Label | ConfigMap Name | Status |
|----------------|----------------|--------|
| `mycc_v1` | `mycc_v1` | ✅ Valid |
| `asset-transfer` | `asset-transfer` | ✅ Valid |
| `basic_1.0` | `basic_1.0` | ✅ Valid |
| `mycc_v1` | `mycc-v1-config` | ❌ Invalid - doesn't match |
| `mycc_v1` | `mycc` | ❌ Invalid - doesn't match |

## Environment Variable Priority

Environment variables are loaded in this order (later values override earlier ones):

1. **Fabric Core Variables** - Built-in variables like `CORE_CHAINCODE_ID_NAME`, `CORE_PEER_ADDRESS`, etc.
2. **ConfigMap Variables** - Variables from the ConfigMap matching the chaincode label

## Complete Example

### Step 1: Create ConfigMap First

```bash
kubectl create configmap asset_transfer_v1 \
  --from-literal=LOG_LEVEL=info \
  --from-literal=MAX_CONNECTIONS=100 \
  --from-literal=CACHE_TTL=3600 \
  -n hyperledger
```

### Step 2: Package Chaincode

```bash
peer lifecycle chaincode package asset-transfer.tar.gz \
  --path ./asset-transfer-chaincode \
  --lang k8s \
  --label asset_transfer_v1
```

### Step 3: Install and Deploy

```bash
# Install
peer lifecycle chaincode install asset-transfer.tar.gz

# Get package ID
peer lifecycle chaincode queryinstalled

# Approve (replace PACKAGE_ID with actual value)
peer lifecycle chaincode approveformyorg \
  --channelID mychannel \
  --name asset-transfer \
  --version 1.0 \
  --package-id asset_transfer_v1:abc123...

# Commit
peer lifecycle chaincode commit \
  --channelID mychannel \
  --name asset-transfer \
  --version 1.0
```

The chaincode pod will now have `LOG_LEVEL`, `MAX_CONNECTIONS`, and `CACHE_TTL` environment variables available.

## Updating Configuration

To update environment variables without redeploying chaincode:

### Option 1: Edit ConfigMap Directly

```bash
kubectl edit configmap mycc_v1 -n hyperledger
```

### Option 2: Apply Updated YAML

```bash
kubectl apply -f mycc-configmap.yaml
```

### Restart Chaincode Pod

After updating the ConfigMap, restart the chaincode pod to pick up changes:

```bash
# Find the pod
kubectl get pods -l fabric-builder-k8s-cclabel=mycc_v1 -n hyperledger

# Delete it (peer will recreate it automatically)
kubectl delete pod <pod-name> -n hyperledger
```

Or delete by label:

```bash
kubectl delete pod -l fabric-builder-k8s-cclabel=mycc_v1 -n hyperledger
```

## Multiple Versions

You can have different configurations for different versions of the same chaincode:

```bash
# Version 1
kubectl create configmap mycc_v1 \
  --from-literal=FEATURE_X=disabled \
  -n hyperledger

# Version 2
kubectl create configmap mycc_v2 \
  --from-literal=FEATURE_X=enabled \
  -n hyperledger
```

Each version will use its own ConfigMap based on the label.


### Wrong ConfigMap Name

Ensure the ConfigMap name exactly matches the chaincode label. Check your package label:

```bash
peer lifecycle chaincode queryinstalled
```

The label is shown in the package ID: `label:hash`

## Best Practices

1. **Create ConfigMap Before Deployment**: Create the ConfigMap before installing the chaincode to ensure it's available immediately
2. **Use Descriptive Labels**: Use clear, version-specific labels like `mycc_v1`, `mycc_v2` instead of generic names
3. **Document Variables**: Add comments in your ConfigMap YAML to document what each variable does
4. **Version Control**: Store ConfigMap YAML files in version control alongside your chaincode
5. **Environment-Specific ConfigMaps**: Use different ConfigMaps for dev, staging, and production environments