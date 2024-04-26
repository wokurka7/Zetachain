# Build Stage
FROM golang:1.20-alpine3.18 AS builder

ENV GOPATH /go
ENV GOOS=linux
ENV CGO_ENABLED=1

# Install build dependencies
RUN apk --no-cache add git make build-base jq openssh libusb-dev linux-headers bash curl python3 py3-pip

# Set the working directory
WORKDIR /go/delivery/zeta-node

# Copy module files and download dependencies
COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy the rest of the source code and build the application
COPY . .

RUN EXPECTED_MAJOR_VERSION=$(grep 'const releaseVersion = ' app/setup_handlers.go | awk -F'"' '{print $2}') && \
    make install VERSION="${EXPECTED_MAJOR_VERSION}" && \
    echo -n "${EXPECTED_MAJOR_VERSION}" > /go/delivery/zeta-node/EXPECTED_MAJOR_VERSION

# Run Stage
FROM alpine:3.18

ENV COSMOVISOR_CHECKSUM="626dfc58c266b85f84a7ed8e2fe0e2346c15be98cfb9f9b88576ba899ed78cdc"
ENV COSMOVISOR_VERSION="v1.5.0"
# Copy Start Script Helpers
COPY contrib/docker-scripts/* /scripts/

# Install runtime dependencies
RUN apk --no-cache add git jq bash curl nano vim tmux python3 libusb-dev linux-headers make build-base bind-tools psmisc coreutils wget py3-pip qemu-img qemu-system-x86_64 && \
    pip install requests && \
    chmod a+x -R /scripts && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.31-r0/glibc-2.31-r0.apk && \
    apk add --force-overwrite --allow-untrusted glibc-2.31-r0.apk && \
    curl https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz > /tmp/google-cloud-sdk.tar.gz && \
    mkdir -p /usr/local/gcloud && \
    tar -C /usr/local/gcloud -xvf /tmp/google-cloud-sdk.tar.gz && \
    /usr/local/gcloud/google-cloud-sdk/install.sh --quiet && \
    ln -s /usr/local/gcloud/google-cloud-sdk/bin/gcloud /usr/bin/gcloud && \
    python /scripts/install_cosmovisor.py

# Copy the binaries from the build stage
COPY --from=builder /go/bin/zetaclientd /usr/local/bin/zetaclientd
COPY --from=builder /go/bin/zetacored /usr/local/bin/zetacored
COPY --from=builder /go/delivery/zeta-node/EXPECTED_MAJOR_VERSION /scripts/EXPECTED_MAJOR_VERSION

# Set the working directory
WORKDIR /usr/local/bin

# Set the default shell
ENV SHELL /bin/bash

EXPOSE 26656
EXPOSE 1317
EXPOSE 8545
EXPOSE 8546
EXPOSE 9090
EXPOSE 26657
EXPOSE 9091