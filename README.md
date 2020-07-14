# bearer-auth

<div align="center">

![alt text](.res/logo.png "bearer-auth")

[![Build Status](https://cloud.drone.io/api/badges/Sykkro/bearer-auth/status.svg)](https://cloud.drone.io/Sykkro/bearer-auth)
[![buildx](https://github.com/Sykkro/bearer-auth/workflows/buildx/badge.svg)](https://github.com/Sykkro/bearer-auth/actions?query=workflow%3Abuildx)
[![Go Report Card](https://goreportcard.com/badge/github.com/sykkro/bearer-auth)](https://goreportcard.com/report/github.com/sykkro/bearer-auth)
[![GitHub All Releases](https://img.shields.io/github/downloads/sykkro/bearer-auth/total)](https://github.com/Sykkro/bearer-auth/releases)

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![Docker Pulls](https://img.shields.io/docker/pulls/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![](https://img.shields.io/microbadger/layers/sykkro/bearer-auth)](https://microbadger.com/images/sykkro/bearer-auth)

<!--
[![Docker Automated](https://img.shields.io/docker/cloud/automated/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![Docker Build](https://img.shields.io/docker/cloud/build/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)-->

</div>

# About 
Traefik forward-auth middleware to impersonate kubernetes user/service accounts for authenticated users.

> *WARNING* this is currently in POC phase, so expect things to break or not work at all as you'd expect

Conceptually, this is a simple HTTP server that runs on port 8000 by default and processes user ids via HTTP headers.
The end purpose is to serve as an enrichement middleware chained after `thomseddon/traefik-forward-auth`, mapping `X-Forwarded-User` users to configured impersonations, by sending a (pre-configured, pod-mounted) impersonator bearer token with `Authentication: Bearer <token>` and impersonated account name with `Impersonate-User` in the response HTTP headers.

Please refer to [this file](test/config_reference.yaml) for configuration guidelines.

*Please note:* This is a WIP project and I have zero knowledge in go, being these my very first lines of code in this language.
Feel free to share suggestions or contribute with improvements/some refactoring. 🛠

# Running

Launch with:
```
./main --config=test/test_config.yaml
```

Test user impersonation with:
```
curl -i http://localhost:8080 -H 'X-Forwarded-User: admin@test.example'
```

# Building

## For local run
```bash
go build main.go
```

## Docker image
```bash
# multi-platform build with buildx
docker buildx build \
--platform=linux/amd64,linux/arm/v5,linux/arm/v6,linux/arm/v7,linux/arm64 \
--output "type=image,push=false" \
-t sykkro/bearer-auth:latest .

```