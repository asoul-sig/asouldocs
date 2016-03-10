FROM alpine:latest
RUN apk add --update git ca-certificates && \
    rm -rf /var/cache/apk/*
ADD custom /custom
ADD peach /bin/peach
ADD public /public
ADD templates /templates
EXPOSE 5556
ENTRYPOINT ["/bin/peach"]
