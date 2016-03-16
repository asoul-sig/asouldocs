#
# Go Dockerfile
#
# https://github.com/peachdocs/peach
#

# Pull base image.
FROM alpine:latest


RUN apk update && apk add curl git mercurial bzr go && rm -rf /var/cache/apk/*



# Set environment variables.
ENV GOROOT /usr/lib/go
ENV GOPATH /gopath
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH


#Copy sourcecode for build
COPY .  ${GOPATH}/src/github.com/peachdocs/peach

# Build peach for golang
RUN cd ${GOPATH}/src/github.com/peachdocs/peach && go get
RUN cd ${GOPATH}/src/github.com/peachdocs/peach && go build   .


# Define working directory.
WORKDIR /app

RUN rm -rf /app/* && cp -r -f  ${GOPATH}/src/github.com/peachdocs/peach/* /app

RUN rm -rf ${GOPATH}
# Define mountable directories.
VOLUME ["/app/custom"]



# Define default command.
CMD ["/app/peach"]


# Expose ports.
EXPOSE 5556