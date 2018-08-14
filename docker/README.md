# Docker for Peach

## Usage

To keep your data out of Docker container, we do a volume (`/var/peach` -> `/data/peach`) here, and you can change it based on your situation.

```
# Pull image from Docker Hub.
$ docker pull peachdocs/peach

# Create local directory for volume.
$ mkdir -p /var/peach

# Use `docker run` for the first time. 
# Peach will complain about missing custom app.ini, leave it there for the moment,
# you will restart the container after finished the configuration part.
$ docker run --name=peach -p 5555:5555 -v /var/peach:/data/peach peachdocs/peach

# Use `docker start` if you have stopped it.
$ docker start peach
```

Files will be store in local path `/var/peach` in my case.

Directory `/var/peach` keeps Git repositories and Gogs data:

    /var/peach
    |-- custom
    |-- data
    |-- log

### Volume with data container

If you're more comfortable with mounting data to a data container, the commands you execute at the first time will look like as follows:

```
# Create data container
docker run --name=peach-data --entrypoint /bin/true peachdocs/peach

# Use `docker run` for the first time.
docker run --name=peach --volumes-from peach-data -p 5555:5555 peachdocs/peach
```

#### Using Docker 1.9 Volume command

```
# Create docker volume.
$ docker volume create --name peach-data

# Use `docker run` for the first time.
$ docker run --name=peach -p 5555:5555 -v peach-data:/data/peach peachdocs/peach
```

## Upgrade

:exclamation::exclamation::exclamation:<span style="color: red">**Make sure you have volumed data to somewhere outside Docker container**</span>:exclamation::exclamation::exclamation:

Steps to upgrade Peach with Docker:

- `docker pull peachdocs/peach`
- `docker stop peach`
- `docker rm peach`
- Finally, create container as the first time and don't forget to do same volume and port mapping.

## Known Issues

- The docker container can not currently be build on Raspberry 1 (armv6l) as our base image `alpine` does not have a `go` package available for this platform.
