# Chaincode job

The k8s builder runs chaincode images using a long running [Kubernetes job](https://kubernetes.io/docs/concepts/workloads/controllers/job/). Using jobs instead of bare pods [enables Kubernetes to clean up chaincode pods automatically](https://kubernetes.io/docs/concepts/workloads/controllers/ttlafterfinished/).

The k8s builder uses labels and annotations to help identify the Kubernetes objects it creates.

## Labels

Kubernetes objects created by the k8s builder have the following labels.

app.kubernetes.io/name[^1]

: The name of the application, `hyperledger-fabric`

app.kubernetes.io/component[^1]

: The application component, `chaincode`

app.kubernetes.io/created-by[^1]

: The tool that created the object, `fabric-builder-k8s`

app.kubernetes.io/managed-by[^1]

: The tool used to manage the application, `fabric-builder-k8s`

fabric-builder-k8s-cclabel

: The chaincode label, e.g. `mycc`

fabric-builder-k8s-cchash

: Base32 encoded chaincode hash, e.g. `U7FELJ6MQXY5RHEQLN3VSIBWD3IITI3E4EVJW3KVXJ24SZO522UQ`

    The chaincode hash is base32 encoded so that it fits in the maximum number of characters allowed for a Kubernetes label value. For example, if you have the chaincode package ID, use the following commands to base32 encode the chaincode hash.

    ```shell
    echo $PACKAGE_ID | cut -d':' -f2 | xxd -r -p | base32 | tr -d '='
    ```

[^1]:
    Kubernetes defines [recommended labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/) to describe applications and instances of applications.

## Annotations

Kubernetes objects created by the k8s builder have the following annotations.

fabric-builder-k8s-ccid

: The full chaincode package ID, e.g. `mycc:a7ca45a7cc85f1d89c905b775920361ed089a364e12a9b6d55ba75c965ddd6a9`

fabric-builder-k8s-mspid

: The membership service provider ID, e.g. `DigiBank`

fabric-builder-k8s-peeraddress

: The peer address, e.g. `peer0.digibank.example.com`

fabric-builder-k8s-peerid

: The peer ID, e.g. `peer0`

