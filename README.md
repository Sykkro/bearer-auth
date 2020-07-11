# bearer-auth

[![Build Status](https://cloud.drone.io/api/badges/Sykkro/bearer-auth/status.svg)](https://cloud.drone.io/Sykkro/bearer-auth)
[![Docker Automated](https://img.shields.io/docker/automated/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![Docker Build](https://img.shields.io/docker/build/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)
[![Docker Pulls](https://img.shields.io/docker/pulls/sykkro/bearer-auth)](https://hub.docker.com/repository/docker/sykkro/bearer-auth)

## About 
Traefik forward-auth middleware to inject bearer tokens for authenticated users.

Conceptually, this is a simple HTTP server that runs on port 8000 by default and processes user tokens via HTTP headers.
The end purpose is to serve as an enrichement middleware chained after `thomseddon/traefik-forward-auth`, mapping `X-Forwarded-User` users to bearer tokens to be injected via `Authentication: Bearer <token>` headers.

This is currently in POC phase, and the only thing this is doing right now is intercepting authenticated users and logging some stuff,
basically acting as a logging MITM proxy for authenticated users.

*Please note:* This is a WIP project and I have zero knowledge in go, being these my very first lines of code in this language.
Feel free to share suggestions or contribute with improvements/some refactoring. 🛠

## Building

``` bash
# AMD64
docker build -t sykkro/bearer-auth:latest -f Dockerfile.amd64 .
# ARM
docker build -t sykkro/bearer-auth:arm-latest -f Dockerfile.arm .
# ARM64
docker build -t sykkro/bearer-auth:arm64-latest -f Dockerfile.arm64 .
```