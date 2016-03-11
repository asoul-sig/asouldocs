#
# Go Dockerfile
#
# https://github.com/peachdocs/peach
#

# Pull base image.
FROM dockerclub/ubuntu

# Install Go
RUN \
  mkdir -p /goroot && \
  curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

# Set environment variables.
ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH


#Copy sourcecode for build
COPY .  ${GOPATH}/src/github.com/peachdocs/peach

# Build peach for golang
RUN cd ${GOPATH}/src/github.com/peachdocs/peach && go get
RUN cd ${GOPATH}/src/github.com/peachdocs/peach && go build -a -tags "netgo static_build" -installsuffix netgo  .


# Define working directory.
WORKDIR /app
COPY .  /app

# Define mountable directories.
VOLUME ["/app/custom"]



# Define default command.
CMD ["/app/peach"]


# Expose ports.
EXPOSE 5556