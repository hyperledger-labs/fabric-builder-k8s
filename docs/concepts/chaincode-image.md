# Chaincode image

Unlike the traditional built-in chaincode language support for Go, Java, and Node.js, the k8s builder *does not* build a chaincode Docker image using Docker-in-Docker.
Instead, a chaincode Docker image must be built and published before it can be used with the k8s builder.

The chaincode will have access to the following environment variables:

- CORE_CHAINCODE_ID_NAME
- CORE_PEER_ADDRESS
- CORE_PEER_TLS_ENABLED
- CORE_PEER_TLS_ROOTCERT_FILE
- CORE_TLS_CLIENT_KEY_PATH
- CORE_TLS_CLIENT_CERT_PATH
- CORE_TLS_CLIENT_KEY_FILE
- CORE_TLS_CLIENT_CERT_FILE
- CORE_PEER_LOCALMSPID

See the [sample contracts for Go, Java, and Node.js](https://github.com/hyperledger-labs/fabric-builder-k8s/tree/main/samples) for basic docker images which will work with the k8s builder.
