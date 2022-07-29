# SPDX-License-Identifier: Apache-2.0

ARG GO_VER=1.17.5
ARG ALPINE_VER=3.14

FROM golang:${GO_VER}-alpine${ALPINE_VER} as build

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

FROM golang:${GO_VER}-alpine${ALPINE_VER}

LABEL org.opencontainers.image.title "Sample Go Contract"
LABEL org.opencontainers.image.description "Sample Hyperledger Fabric Go contract for Kubernetes chaincode builder"
LABEL org.opencontainers.image.source "https://github.com/hyperledger-labs/fabric-builder-k8s/samples/go-contract"

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=build /go/bin/go-contract /usr/bin/go-contract

WORKDIR /var/hyperledger/go-contract
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "exec /usr/bin/go-contract -peer.address=$CORE_PEER_ADDRESS"]
