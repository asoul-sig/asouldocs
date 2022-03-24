#!/bin/sh
set -x
set -e

# Set temp environment vars
export GOPATH=/tmp/go
export PATH=${PATH}:${GOPATH}/bin

# Install build deps
apk --no-cache --no-progress add --virtual build-deps go gcc musl-dev

# Init go environment to build
mkdir -p ${GOPATH}/src/github.com/asoul-go/
ln -s /app/asouldocs/ ${GOPATH}/src/github.com/asoul-go/asouldocs
cd ${GOPATH}/src/github.com/asoul-go/asouldocs
go get -v
mv ${GOPATH}/bin/asouldocs .

# Cleanup GOPATH
rm -r $GOPATH

# Remove build deps
apk --no-progress del build-deps

# Create user
adduser -H -D -g 'ASoulDocs' asouldocs -h /data/asouldocs -s /bin/bash && passwd -u asouldocs
