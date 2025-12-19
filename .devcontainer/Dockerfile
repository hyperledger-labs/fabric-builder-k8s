FROM mcr.microsoft.com/devcontainers/go:1.25-bookworm@sha256:ef7d7fe9f8eb237c2d6b8403bc0affceca3c349a041fb7714dfc554d9a97adca

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
  && apt-get -y install --no-install-recommends \
  xxd
