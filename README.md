# fabric-builder-k8s

Kubernetes [external chaincode builder](https://hyperledger-fabric.readthedocs.io/en/latest/cc_launcher.html) for Hyperledger Fabric.

Advantages:

🚀 Chaincode runs _immediately_ on channel commit.

✨ Avoids the complexity and administrative burdens associated with Chaincode-as-a-Service.

🔥 Pre-published chaincode images avoid code-compilation errors at deployment time.

🏗️ Pre-published chaincode images encourage modern, industry accepted CI/CD best practices.

🛡️ Pre-published chaincode images remove any and all dependencies on a root-level _docker daemon_.

🕵️ Pre-published chaincode images provide traceability and change management features (e.g. Git commit hash as image tag)
