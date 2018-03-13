.PHONY: build run

all: build run

build:
	go install -v
	cp '$(GOPATH)/bin/peach' .

run:
	./peach web

bindata:
	go-bindata -o=pkg/bindata/bindata.go -ignore="\\.DS_Store|README|config.codekit|.less" -pkg=bindata templates/... conf/... public/...