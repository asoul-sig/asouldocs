
# Docker for peach

Visit [Docker Hub](https://hub.docker.com/r/peachdocs/peach/) see all available tags.

## Usage

To keep your data out of Docker container, we do a volume (`/app/custom` -> `/custom`) here, and you can change it based on your situation.

```
# Pull image from Docker Hub.
$ docker pull peachdocs/peach

# Create local directory for volume.
$ mkdir -p /custom

# Use `make dockerrun` or `docker run` for the first time.

$ make dockerrun PERACH_PORT={YOUR HOST PORT} PERACH_CUSTOM_PATH=/custom

or 

$docker run -ti  -d  -p {YOUR HOST PORT}:5556 --restart=always  --entrypoint="/bin/bash" --name peach -v /custom:/app/custom  peachdocs/peach /app/peach web

```

## About make command

### local build peach

$ make all


### local test run peach application

$ make testrun


### make a peach docker image

$ make image

### push a peach docker image to docker hub

$ make release