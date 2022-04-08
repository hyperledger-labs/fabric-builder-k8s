ARG GO_VER=1.17.5
ARG ALPINE_VER=3.14
ARG HLF_VERSION=2.4

FROM golang:${GO_VER}-alpine${ALPINE_VER} as build
RUN apk add --no-cache \
	bash \
	binutils-gold \
	gcc \
	git \
	make \
	musl-dev

ADD . $GOPATH/src/github.com/hyperledgendary/fabric-builder-k8s
WORKDIR $GOPATH/src/github.com/hyperledgendary/fabric-builder-k8s

RUN go install ./cmd/...

RUN ls -al $GOPATH/bin

FROM hyperledger/fabric-peer:${HLF_VERSION} as core

RUN apk add --no-cache \
	wget;

RUN wget https://github.com/mikefarah/yq/releases/download/v4.23.1/yq_linux_amd64 -O /usr/bin/yq && chmod +x /usr/bin/yq

RUN yq 'del(.vm.endpoint) | .chaincode.externalBuilders += { "name": "k8s_builder", "path": "/opt/hyperledger/k8s_builder", "propagateEnvironment": [ "KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT" ] }' ${FABRIC_CFG_PATH}/core.yaml > core.yaml

FROM hyperledger/fabric-peer:${HLF_VERSION}

COPY --from=core core.yaml ${FABRIC_CFG_PATH}
COPY --from=build /go/bin/ /opt/hyperledger/k8s_builder/bin/
