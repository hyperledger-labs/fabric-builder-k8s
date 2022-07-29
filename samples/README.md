# Sample contracts

The main purpose of these samples is to demonstrate basic `Dockerfile`s for deploying Go, Java, and Node.js chaincode with the k8s builder.

The samples can be built with:

```shell
docker build . -t sample-contract
```

You will need a digest to create a chaincode package for the k8s builder, which is only created when the docker image is published to a registry.
For example, to publish to a local registry:

```shell
docker tag sample-contract localhost:5000/sample-contract
docker push localhost:5000/sample-contract
```
