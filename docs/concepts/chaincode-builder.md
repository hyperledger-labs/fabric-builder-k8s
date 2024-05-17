# Chaincode builder

From version 2.0, Hyperledger Fabric supports External Builders and Launchers to manage the process of building and launching chaincode, rather than being limited to the peer's built in Docker based build and launch process.

External Builders and Launchers are a collection of executables — `detect`, `build`, `release` and `run` — that the peer calls in order to build, launch, and discover chaincode. The k8s builder executables do the following.

| Executable  | Description                                                                                                       |
| ----------- | ----------------------------------------------------------------------------------------------------------------- |
| detect      | Detects chaincode packages with a type of `k8s`                                                                   |
| build       | No-op (a chaincode image must be built and published to a container registry before it can be deployed to Fabric) |
| release     | No-op                                                                                                             |
| run         | Starts a Kubernetes pod using the chaincode image identified by an immutable digest                               |

For more information about Fabric builders, see the [External Builders and Launchers](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) documentation.
