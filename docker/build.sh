#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin

# Install build deps
apk --no-cache --no-progress add --virtual build-deps go gcc musl-dev

# Init go environment to build
mkdir -p ${GOPATH}/src/github.com/peachdocs/
ln -s /app/peach/ ${GOPATH}/src/github.com/peachdocs/peach
cd ${GOPATH}/src/github.com/peachdocs/peach
go get -v
mv ${GOPATH}/bin/peach .

# Cleanup GOPATH
rm -r $GOPATH

# Remove build deps
apk --no-progress del build-deps

# Create user
adduser -H -D -g 'Peach Docs' peach -h /data/peach -s /bin/bash && passwd -u peach
