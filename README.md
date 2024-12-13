# fabric-builder-k8s

[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/9817/badge)](https://www.bestpractices.dev/projects/9817)

The Kubernetes [external chaincode builder](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) for Hyperledger Fabric (k8s builder) is an alternative to Fabric's legacy built in Docker chaincode builder, which does not work in a Kubernetes deployment, and the preconfigured chaincode-as-a-service builder, which is more suited to chaincode development and test.

For more information, including how to deploy your first chaincode with the k8s builder, see the [k8s builder documentation](https://labs.hyperledger.org/fabric-builder-k8s/).

To find out how to report issues, suggest enhancements and contribute to the k8s builder project, see the [contributing guide](CONTRIBUTING.md).

## Overview

With the k8s builder, the Fabric administrator is responsible for preparing a [chaincode image](https://labs.hyperledger.org/fabric-builder-k8s/concepts/chaincode-image/), publishing to a container registry, and preparing a [chaincode package](https://labs.hyperledger.org/fabric-builder-k8s/concepts/chaincode-package/) with coordinates of the contract's immutable image digest.
When Fabric detects the installation of a `type=k8s` contract, the builder assumes full ownership of the lifecycle of pods, containers, and network linkages necessary to communicate securely with the peer.

Advantages:

🚀 Chaincode runs _immediately_ on channel commit.

✨ Avoids the complexity and administrative burdens associated with Chaincode-as-a-Service.

🔥 Pre-published chaincode images avoid code-compilation errors at deployment time.

🏗️ Pre-published chaincode images encourage modern, industry accepted CI/CD best practices.

🛡️ Pre-published chaincode images remove any and all dependencies on a root-level _docker daemon_.

🕵️ Pre-published chaincode images provide traceability and change management features (e.g. Git commit hash as image tag)
