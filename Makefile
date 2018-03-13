.PHONY: build run

all: build run

build:
	go install -v
	cp '$(GOPATH)/bin/peach' .

run:
	./peach web