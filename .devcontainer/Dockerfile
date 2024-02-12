FROM mcr.microsoft.com/devcontainers/go:1-1-bookworm

# Install mongosh
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get -y install --no-install-recommends gnupg && \
    wget -qO- https://www.mongodb.org/static/pgp/server-7.0.asc | tee /etc/apt/trusted.gpg.d/server-7.0.asc && \
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-7.0.list && \
    apt-get update && \
    apt-get -y install --no-install-recommends mongodb-mongosh

# Install k6
ARG TARGETARCH
ARG K6_VERSION=0.49.0

RUN wget -q -O k6.tar.gz "https://github.com/grafana/k6/releases/download/v${K6_VERSION}/k6-v${K6_VERSION}-linux-${TARGETARCH}.tar.gz" \
  && tar xzf k6.tar.gz \
  && mv k6-v${K6_VERSION}-linux-${TARGETARCH}/k6 /usr/bin/ \
  && chmod +x /usr/bin/k6 \
  && rm -rf k6-v${K6_VERSION}-linux-${TARGETARCH} \
  && rm k6.tar.gz
