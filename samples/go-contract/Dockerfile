# SPDX-License-Identifier: Apache-2.0

ARG GO_VER=1.26.0
ARG ALPINE_VER=3.21
ARG IMAGE_DIGEST=sha256:7a3e50096189ad57c9f9f865e7e4aa8585ed1585248513dc5cda498e2f41812c

FROM golang:${GO_VER}-alpine${ALPINE_VER}@${IMAGE_DIGEST} AS build

RUN apk add --no-cache \
	bash \
	binutils-gold \
  dumb-init \
	gcc \
	git \
	make \
	musl-dev

ADD . $GOPATH/src/github.com/hyperledger-labs/fabric-builder-k8s/samples/go-contract
WORKDIR $GOPATH/src/github.com/hyperledger-labs/fabric-builder-k8s/samples/go-contract

RUN go install ./...

FROM golang:${GO_VER}-alpine${ALPINE_VER}@${IMAGE_DIGEST}

LABEL org.opencontainers.image.title="Sample Go Contract"
LABEL org.opencontainers.image.description="Sample Hyperledger Fabric Go contract for Kubernetes chaincode builder"
LABEL org.opencontainers.image.source="https://github.com/hyperledger-labs/fabric-builder-k8s/samples/go-contract"

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=build /go/bin/go-contract /usr/bin/go-contract

WORKDIR /var/hyperledger/go-contract
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "exec /usr/bin/go-contract -peer.address=$CORE_PEER_ADDRESS"]
