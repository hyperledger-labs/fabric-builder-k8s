# SPDX-License-Identifier: Apache-2.0

ARG JAVA_VERSION=25

FROM eclipse-temurin:${JAVA_VERSION}-jdk-alpine AS build

RUN apk add --no-cache \
    dumb-init wget unzip

ARG GRADLE_VERSION=9.7.1
ARG GRADLE_SHA256=acd53f1edaf02f1a8ff99879f8a34b302661a057d9b063ae9e35b552f804d20a

RUN wget -q "https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip" -O /tmp/gradle.zip && \
    echo "${GRADLE_SHA256}  /tmp/gradle.zip" | sha256sum -c - && \
    unzip -qq /tmp/gradle.zip -d /opt && \
    mv /opt/gradle-${GRADLE_VERSION} /opt/gradle && \
    rm -f /tmp/gradle.zip

ENV PATH=/opt/gradle/bin:${PATH}

WORKDIR /usr/src/app

COPY . .

RUN gradle jar

FROM eclipse-temurin:${JAVA_VERSION}-jre-alpine

LABEL org.opencontainers.image.title="Sample Java Contract"
LABEL org.opencontainers.image.description="Sample Hyperledger Fabric Java contract for Kubernetes chaincode builder"
LABEL org.opencontainers.image.source="https://github.com/hyperledger-labs/fabric-builder-k8s/samples/java-contract"

WORKDIR /var/hyperledger/java-contract

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=build /usr/src/app/build/libs/sample-contract.jar ./sample-contract.jar

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "exec /opt/java/openjdk/bin/java -jar ./sample-contract.jar --peer.address=$CORE_PEER_ADDRESS"]
