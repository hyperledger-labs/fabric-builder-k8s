# Developing and debugging chaincode

Publishing a chaincode Docker image and using the image digest to deploy the chaincode works well with the [Fabric chaincode lifecycle](https://hyperledger-fabric.readthedocs.io/en/latest/chaincode_lifecycle.html), however it is not as convenient while developing and debugging chaincode.

This tutorial describes how to debug chaincode using the chaincode-as-a-service (CCAAS) builder before you are ready to deploy it in a production environment with the k8s builder.

With CCAAS you still have to go through the chaincode lifecycle once to tell Fabric where the chaincode is running but after that you can stop, update, restart, and debug the chaincode all without needing to repeat the chaincode lifecycle steps as you would do normally.

Like the k8s builder, CCAAS packages do not contain any chaincode source code.
Unlike the k8s builder, CCAAS packages do not uniquely identify any specific chaincode: they only contain details a Fabric peer needs to connect to a chaincode instance running in server mode.
This is unlikely to be acceptable in a production environment but it is ideal in a development environment.

First create a directory to download all the required files and run the demo.

```shell
mkdir hlf-debug-demo
cd hlf-debug-demo
```

Now follow the steps below to debug your first smart contract using the CCAAS builder!

## Setup the nano test network

In this tutorial, we'll use the [fabric-samples](https://github.com/hyperledger/fabric-samples/) nano test network because it's ideally suited to a light weight development environment.
There is also a [test-network CCAAS tutorial](https://github.com/hyperledger/fabric-samples/blob/main/test-network/CHAINCODE_AS_A_SERVICE_TUTORIAL.md) if you would prefer to use the Docker based test network.

Start by downloading the sample nano test network (fabric-samples isn't tagged so we'll use a known good commit).

```shell
export FABRIC_SAMPLES_COMMIT=0db64487e5e89a81d68e6871af3f0907c67e7d75
curl -sSL "https://github.com/hyperledger/fabric-samples/archive/${FABRIC_SAMPLES_COMMIT}.tar.gz" | tar -xzf - --strip-components=1 fabric-samples-${FABRIC_SAMPLES_COMMIT}/test-network-nano-bash
```

You will also need to install the Fabric binaries, which include the CCAAS builder and default configuration files.

```shell
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh binary
```

The default `core.yaml` file provided with the Fabric binaries needs to be updated with the correct CCAAS builder location, since it assumes that the peer will be running in a docker container. Either update the `ccaas_builder` in `externalBuilders` with the absolute path to the `builders/ccaas` directory using your preferred editor, or use the following [yq](https://github.com/mikefarah/yq) command to update the path.

```shell
yq -i '( .chaincode.externalBuilders[] | select(.name == "ccaas_builder") | .path ) = strenv(PWD) + "/builders/ccaas"' config/core.yaml
```

Now start the nano test network!

```shell
cd test-network-nano-bash
./network.sh start
```

## Deploy a CCAAS chaincode package

Open a new shell, change to the `hlf-debug-demo/test-network-nano-bash` directory, and check that the nano test network is running.

```shell
. ./peer1admin.sh
peer channel list
```

Sourcing the `peer1admin.sh` script sets up the environment for running `peer` commands.
You'll need to repeat this step in any new shell you want to run `peer` commands in, for example to deploy chaincode or run transactions.

Before deploying a CCAAS chaincode package, we need to create a suitable chaincode package.
Start by creating a `connection.json` file.
You'll need to use the same `address` when starting the chaincode later. 

```shell
cat << CONNECTIONJSON_EOF > connection.json
{
  "address": "127.0.0.1:9999",
  "dial_timeout": "10s",
  "tls_required": false
}
CONNECTIONJSON_EOF
```

Next, create a `metadata.json` file for the chaincode package.
The CCAAS builder provided with Fabric will detect the type of `ccaas`.

```shell
cat << METADATAJSON_EOF > metadata.json
{
    "type": "ccaas",
    "label": "dev-contract"
}
METADATAJSON_EOF
```

Create the chaincode package archive.

```shell
tar -czf code.tar.gz connection.json
tar -czf dev-contract.tgz metadata.json code.tar.gz
```

Now follow the Fabric chaincode lifecycle to deploy the CCAAS chaincode.
Start by installing the chaincode package you just created.

```shell
peer lifecycle chaincode install dev-contract.tgz
```

Set the CHAINCODE_ID environment variable for use in subsequent commands.

```shell
export CHAINCODE_ID=$(peer lifecycle chaincode calculatepackageid dev-contract.tgz) && echo $CHAINCODE_ID
```

Approve and commit the chaincode.

```shell
peer lifecycle chaincode approveformyorg -o 127.0.0.1:6050 --channelID mychannel --name dev-contract --version 1 --package-id $CHAINCODE_ID --sequence 1 --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
peer lifecycle chaincode commit -o 127.0.0.1:6050 --channelID mychannel --name dev-contract --version 1 --sequence 1 --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

## Debugging the chaincode

Normally Fabric peers start chaincode automatically but with CCAAS you are responsible for starting the chaincode in server mode manually, which is why it is so useful for debugging.

Importantly, chaincode can be started in server mode without any code changes in recent Fabric versions, making it easy to move from CCAAS in a development environment to using the k8s builder in a production environment.

Go and Java chaincode will start in server mode if `CHAINCODE_SERVER_ADDRESS` and `CORE_CHAINCODE_ID_NAME` environment variables are configured.
The `CHAINCODE_SERVER_ADDRESS` must match the `address` field from the `connection.json` file in the chaincode package.
The `CORE_CHAINCODE_ID_NAME` can be found using the `peer lifecycle chaincode calculatepackageid` command.

Node.js chaincode can be started in server mode using the `server` command and providing `--chaincode-address` and `--chaincode-id` command line arguments.
The sample Node.js contract includes a `debug` script in `package.json` which uses the same `CHAINCODE_SERVER_ADDRESS` and `CORE_CHAINCODE_ID_NAME` environment variables for these command line arguments as Go and Java chaincode.

For example, use the following commands to export the required environment variables.

```shell
export CHAINCODE_SERVER_ADDRESS=127.0.0.1:9999
export CORE_CHAINCODE_ID_NAME=$(peer lifecycle chaincode calculatepackageid dev-contract.tgz)
```

The following commands can then be used to start each sample in server mode without a debugger attached.

In the `samples/go-contract` directory:

```shell
CORE_PEER_TLS_ENABLED=false go run main.go
```

In the `samples/java-contract` directory:

```shell
./gradlew jar
CORE_PEER_TLS_ENABLED=false java -jar ./build/libs/sample-contract.jar
```

In the `samples/node-contract` directory:

```shell
npm install
npm run compile
npm run debug
```

Change to the `hlf-debug-demo/test-network-nano-bash` directory in a new shell, and check everything is working using the `GetMetadata` transaction.

```shell
. ./peer1admin.sh
peer chaincode query -C mychannel -n dev-contract -c '{"Args":["org.hyperledger.fabric:GetMetadata"]}'
```

Now it's time to attach a debugger, set breakpoints and start running transactions!
Exactly how you debug chaincode will depend on your preferred debugger.
For example, [VS Code includes a built-in debugger for Node.js](https://code.visualstudio.com/docs/nodejs/nodejs-debugging).

Try adding breakpoints to the transactions defined in the samples and then invoke the `PutValue` transaction using the following command.

```shell
peer chaincode invoke -o 127.0.0.1:6050 -C mychannel -n dev-contract -c '{"Args":["PutValue","asset1","green"]}' --waitForEvent --tls --cafile "${PWD}"/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls/ca.crt
```

Query the `GetValue` transaction using the following command.

```shell
peer chaincode query -C mychannel -n dev-contract -c '{"Args":["GetValue","asset1"]}'
```

## Next steps

Take a look at the [Fabric Full Stack Development Workshop](https://github.com/hyperledger/fabric-samples/blob/main/full-stack-asset-transfer-guide/README.md) for an in-depth introduction to the entire Fabric development process, from smart contract and client application development, to cloud native deployment.
