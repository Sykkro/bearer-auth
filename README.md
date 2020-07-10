# bearer-auth

[![Build Status](https://cloud.drone.io/api/badges/Sykkro/bearer-auth/status.svg)](https://cloud.drone.io/Sykkro/bearer-auth)

## About 
Simple HTTP server that runs on port 8000 by default and processes user tokens via HTTP headers.
The end purpose is to serve as a traefik forward auth middleware, mapping X-Forwarded-User users to bearer tokens.

This is currently in POC phase, and the only thing this is doing right now is intercepting authenticated users and logging some stuff,
basically acting as a MITM proxy for authenticated users through `thomseddon/traefik-forward-auth`.

*Please note:* This is a WIP project and I have 0 knowledge in go, being these my very first lines of code in this language.
Feel free to share suggestions or contribute with improvements/some refactoring. 🛠

## Building

```
docker build --build-arg go_opts="CGO_ENABLED=0 GOARCH=amd64" --pull -t $IMAGE:amd64 .
docker push $IMAGE:amd64
```
```
docker build --build-arg go_opts="GOARCH=arm GOARM=6" --pull -t $IMAGE:arm32v6 .
docker push $IMAGE:arm32v6
```