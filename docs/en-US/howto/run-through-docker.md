---
title: Run through Docker
---

Docker images for the _**ASoulDocs**_ server are available both on [Docker Hub](https://hub.docker.com/r/unknwon/asouldocs) and [GitHub Container Registry](https://github.com/asoul-sig/asouldocs/pkgs/container/asouldocs).

The `latest` tag represents the latest build from the [`main` branch](https://github.com/asoul-sig/asouldocs).

## Caveats

The `HTTP_ADDR` should be changed to listen on the Docker network or all network addresses:

```ini
HTTP_ADDR = 0.0.0.0
```

## Start the container

You need to volume the `custom` directory into the Docker container for it being able to start (`/app/asouldocs/custom` is the path inside the container):

```bash
$ docker run \
    --name=asouldocs \
    -p 15555:5555 \
    -v $(pwd)/custom:/app/asouldocs/custom \
    unknwon/asouldocs
```

You should also volumn the `docs` directory into the Docker container directly if you are not using a [remote Git address](set-up-documentation.md#target) (`/app/asouldocs/docs` is the path inside the container):

```bash
$ docker run \
    --name=asouldocs \
    -p 15555:5555 \
    -v $(pwd)/custom:/app/asouldocs/custom \
    -v $(pwd)/docs:/app/asouldocs/docs \
    unknwon/asouldocs
```
