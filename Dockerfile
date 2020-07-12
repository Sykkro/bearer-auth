FROM --platform=$BUILDPLATFORM golang:1.13-alpine as base

# Prepare base images with go opts for each target architecture
FROM base AS base-armv5
ENV GO_ARCH_OPS "CGO_ENABLED=0 GOARCH=arm GOARM=5"

FROM base AS base-armv6
ENV GO_ARCH_OPS "CGO_ENABLED=0 GOARCH=arm GOARM=6"

FROM base AS base-armv7
ENV GO_ARCH_OPS "CGO_ENABLED=0 GOARCH=arm GOARM=7"

FROM base AS base-arm64
ENV GO_ARCH_OPS "CGO_ENABLED=0 GOARCH=arm64"

FROM base AS base-amd64
ENV GO_ARCH_OPS "CGO_ENABLED=0 GOARCH=amd64"

# Use the target architecture base image as builder target
ARG TARGETARCH
ARG TARGETVARIANT
FROM base-$TARGETARCH$TARGETVARIANT AS builder

# Setup
RUN mkdir -p /go/src/github.com/sykkro/bearer-auth
WORKDIR /go/src/github.com/sykkro/bearer-auth

# Add libraries
RUN apk add --no-cache git

# Copy & build
ADD . /go/src/github.com/sykkro/bearer-auth/
RUN env ${GO_ARCH_OPS} GOOS=linux GO111MODULE=on go build -a -installsuffix nocgo -o /bearer-auth github.com/sykkro/bearer-auth

# Copy into scratch container
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bearer-auth ./
EXPOSE 8080
ENTRYPOINT ["./bearer-auth"]