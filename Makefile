# this a peach build script
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64

TAG=${TAG:-latest}
COMMIT=`git rev-parse --short HEAD`

PERACH_PORT=${PERACH_PORT:-5556}
PERACH_CUSTOM_PATH=${PERACH_CUSTOM_PATH:-/app/custom}


all: clean init build

clean:
	@rm -rf ./peach

build:
	@go get && go build   .

init:
	@git clone https://github.com/peachdocs/peach.peach custom


image:
	@ echo Building peach image $(TAG)
	@ docker build -t peachdocs/peach:$(TAG) .

release: build image
	@docker push peachdocs/peach:$(TAG) .

test: clean
	@go test -v ./...

testbuild:
	@go build


testrun:
	@ ./peach web


dockerrun:
	@ docker run -ti  -d  -p ${PERACH_PORT}:5556 --restart=always --name peach -v ${PERACH_CUSTOM_PATH}:/app/custom  dockerclub/peach /app/peach web


.PHONY: all build clean  image test release
