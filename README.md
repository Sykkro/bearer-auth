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

<!--
[![Docker Automated](https://img.shields.io/docker/cloud/automated/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![Docker Build](https://img.shields.io/docker/cloud/build/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)-->

</div>

# About 
Traefik forward-auth middleware to inject bearer tokens for authenticated users.

Conceptually, this is a simple HTTP server that runs on port 8000 by default and processes user tokens via HTTP headers.
The end purpose is to serve as an enrichement middleware chained after `thomseddon/traefik-forward-auth`, mapping `X-Forwarded-User` users to bearer tokens to be injected via `Authentication: Bearer <token>` headers.

This is currently in POC phase, and the only thing this is doing right now is intercepting authenticated users and logging some stuff,
basically acting as a logging MITM proxy for authenticated users.

*Please note:* This is a WIP project and I have zero knowledge in go, being these my very first lines of code in this language.
Feel free to share suggestions or contribute with improvements/some refactoring. 🛠

# Building


```bash
# multi-platform build with buildx
docker buildx build \
--platform=linux/amd64,linux/arm/v5,linux/arm/v6,linux/arm/v7,linux/arm64 \
--output "type=image,push=false" \
-t sykkro/bearer-auth:latest .

```