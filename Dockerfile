ARG UBUNTU_VER=24.04
ARG HLF_VERSION=2.5

FROM ubuntu:${UBUNTU_VER} AS build
ARG GO_VER=1.23.0
ENV GOPATH=/go

ENV DEBIAN_FRONTEND="noninteractive"
RUN apt-get update && apt-get install -y -q --no-install-recommends \
    ca-certificates \
    build-essential \
    git \
    gcc \
    curl \
    make

RUN curl -sL https://go.dev/dl/go${GO_VER}.linux-$(dpkg --print-architecture).tar.gz | tar zxf - -C /usr/local
ENV PATH="/usr/local/go/bin:$PATH"

ADD . $GOPATH/src/github.com/hyperledger-labs/fabric-builder-k8s
WORKDIR $GOPATH/src/github.com/hyperledger-labs/fabric-builder-k8s

RUN go install -a -v ./cmd/...

FROM hyperledger/fabric-peer:${HLF_VERSION} AS core

ENV DEBIAN_FRONTEND="noninteractive"
RUN apt-get update && apt-get install -y -q --no-install-recommends \
    ca-certificates \
    wget

RUN wget https://github.com/mikefarah/yq/releases/download/v4.23.1/yq_linux_$(dpkg --print-architecture) -O /usr/bin/yq && chmod +x /usr/bin/yq

RUN yq 'del(.vm.endpoint) | .chaincode.externalBuilders += { "name": "k8s_builder", "path": "/opt/hyperledger/k8s_builder", "propagateEnvironment": [ "CORE_PEER_ID", "KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT" ] }' ${FABRIC_CFG_PATH}/core.yaml > core.yaml

FROM hyperledger/fabric-peer:${HLF_VERSION}

LABEL org.opencontainers.image.title="K8s Hyperledger Fabric Peer"
LABEL org.opencontainers.image.description="Hyperledger Fabric Peer with a preconfigured Kubernetes chaincode builder"
LABEL org.opencontainers.image.source="https://github.com/hyperledger-labs/fabric-builder-k8s"

COPY --from=core core.yaml ${FABRIC_CFG_PATH}
COPY --from=build /go/bin/ /opt/hyperledger/k8s_builder/bin/
