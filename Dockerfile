FROM golang:alpine3.15 AS binarybuilder
RUN apk --no-cache --no-progress add --virtual \
    build-deps \
    build-base \
    git

# Install Task
RUN wget --quiet https://github.com/go-task/task/releases/download/v3.12.0/task_linux_amd64.tar.gz -O task_linux_amd64.tar.gz \
  && sh -c 'echo "803d3c1752da31486cbfb4ddf7d8ba5e0fa8c8ebba8acf227a9cd76ff9901343  task_linux_amd64.tar.gz" | sha256sum -c' \
  && tar -xzf task_linux_amd64.tar.gz \
  && mv task /usr/local/bin/task

WORKDIR /dist
COPY . .
RUN task build

FROM alpine:3.15
RUN echo https://dl-cdn.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
  ca-certificates \
  git

# Install gosu
RUN export url="https://github.com/tianon/gosu/releases/download/1.14/gosu-"; \
  if [ `uname -m` == "aarch64" ]; then \
       wget --quiet ${url}arm64 -O /usr/sbin/gosu \
    && sh -c 'echo "73244a858f5514a927a0f2510d533b4b57169b64d2aa3f9d98d92a7a7df80cea  /usr/sbin/gosu" | sha256sum -c'; \
  elif [ `uname -m` == "armv7l" ]; then \
       wget --quiet ${url}armhf -O /usr/sbin/gosu \
    && sh -c 'echo "abb1489357358b443789571d52b5410258ddaca525ee7ac3ba0dd91d34484589  /usr/sbin/gosu" | sha256sum -c'; \
  else \
       wget --quiet ${url}amd64 -O /usr/sbin/gosu \
    && sh -c 'echo "bd8be776e97ec2b911190a82d9ab3fa6c013ae6d3121eea3d0bfd5c82a0eaf8c  /usr/sbin/gosu" | sha256sum -c'; \
  fi \
  && chmod +x /usr/sbin/gosu

WORKDIR /app/asouldocs/
COPY --from=binarybuilder /dist/asouldocs .

VOLUME ["/app/asouldocs/custom"]
EXPOSE 5555
CMD ["/app/asouldocs/asouldocs", "web"]
