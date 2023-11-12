FROM golang:1.21-alpine3.18 AS binarybuilder
RUN apk --no-cache --no-progress add --virtual \
    build-deps \
    build-base \
    git

# Install Task
RUN wget --quiet https://github.com/go-task/task/releases/download/v3.31.0/task_linux_amd64.tar.gz -O task_linux_amd64.tar.gz \
  && sh -c 'echo "fc707db87c08579066e312b6fd3de4301f7c1c3e48198a86368d063efce2bbab task_linux_amd64.tar.gz" | sha256sum -c' \
  && tar -xzf task_linux_amd64.tar.gz \
  && mv task /usr/local/bin/task

WORKDIR /dist
COPY . .
RUN task build

FROM golang:1.21-alpine3.18
RUN echo https://dl-cdn.alpinelinux.org/alpine/edge/community/ >> /etc/apk/repositories \
  && apk --no-cache --no-progress add \
  ca-certificates \
  git

# Install gosu
RUN export url="https://github.com/tianon/gosu/releases/download/1.17/gosu-"; \
  if [ `uname -m` == "aarch64" ]; then \
       wget --quiet ${url}arm64 -O /usr/sbin/gosu \
    && sh -c 'echo "c3805a85d17f4454c23d7059bcb97e1ec1af272b90126e79ed002342de08389b /usr/sbin/gosu" | sha256sum -c'; \
  elif [ `uname -m` == "armv7l" ]; then \
       wget --quiet ${url}armhf -O /usr/sbin/gosu \
    && sh -c 'echo "e5866286277ff2a2159fb9196fea13e0a59d3f1091ea46ddb985160b94b6841b /usr/sbin/gosu" | sha256sum -c'; \
  else \
       wget --quiet ${url}amd64 -O /usr/sbin/gosu \
    && sh -c 'echo "bbc4136d03ab138b1ad66fa4fc051bafc6cc7ffae632b069a53657279a450de3 /usr/sbin/gosu" | sha256sum -c'; \
  fi \
  && chmod +x /usr/sbin/gosu

WORKDIR /app/asouldocs/
COPY --from=binarybuilder /dist/asouldocs .

VOLUME ["/app/asouldocs/custom"]
EXPOSE 5555
CMD ["/app/asouldocs/asouldocs", "web"]
