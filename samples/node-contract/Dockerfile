# SPDX-License-Identifier: Apache-2.0

ARG NODE_VER=16
ARG ALPINE_VER=3.14

FROM node:${NODE_VER}-alpine${ALPINE_VER} as build

RUN apk add --no-cache \
  dumb-init

FROM node:${NODE_VER}-alpine${ALPINE_VER}

LABEL org.opencontainers.image.title "Sample Node Contract"
LABEL org.opencontainers.image.description "Sample Hyperledger Fabric Node contract for Kubernetes chaincode builder"
LABEL org.opencontainers.image.source "https://github.com/hyperledger-labs/fabric-builder-k8s/samples/node-contract"

COPY --from=build /usr/bin/dumb-init /usr/bin/dumb-init

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run compile

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["sh", "-c", "exec npm start -- --peer.address $CORE_PEER_ADDRESS"]
