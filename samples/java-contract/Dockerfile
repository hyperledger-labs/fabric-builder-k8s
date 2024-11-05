# SPDX-License-Identifier: Apache-2.0

ARG JAVA_VER=11

FROM eclipse-temurin:${JAVA_VER}-jdk-alpine AS build

RUN apk add --no-cache \
  dumb-init

WORKDIR /usr/src/app

COPY ./gradle/wrapper/ ./gradle/wrapper/
COPY gradlew .
RUN ./gradlew --version

COPY . .

RUN ./gradlew jar

FROM eclipse-temurin:${JAVA_VER}-jre-alpine

LABEL org.opencontainers.image.title="Sample Java Contract"
LABEL org.opencontainers.image.description="Sample Hyperledger Fabric Java contract for Kubernetes chaincode builder"
LABEL org.opencontainers.image.source="https://github.com/hyperledger-labs/fabric-builder-k8s/samples/java-contract"

WORKDIR /var/hyperledger/java-contract

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=build /usr/src/app/build/libs/sample-contract.jar ./sample-contract.jar

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "exec /opt/java/openjdk/bin/java -jar ./sample-contract.jar --peer.address=$CORE_PEER_ADDRESS"]
