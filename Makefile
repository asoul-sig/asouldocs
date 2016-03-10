CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
TAG=${TAG:-latest}
COMMIT=`git rev-parse --short HEAD`



all: clean init  build

clean:
	@rm -rf ./peach custom

build:
	@godep go build -a -tags "netgo static_build" -installsuffix netgo  .
#	@cd controller && godep go build -a -tags "netgo static_build" -installsuffix netgo -ldflags "-w -X github.com/dockerclubgroup/shipyard/version.GitCommit=$(COMMIT)" .


init:
	go get github.com/tools/godep
	godep save
	@git clone https://github.com/dockerclubgroup/shipyard.peach custom

image: clean init build
	@ echo Building peach image $(TAG)
	@ docker build -t dockerclub/peach:$(TAG) .

release: build image
	@docker push dockerclub/peach:$(TAG) .

test: clean
	@godep go test -v ./...

testbuild:
	@godep go build -a -tags "netgo static_build" -installsuffix netgo


testrun:
	@ ./peach web

.PHONY: all build clean media image test release
