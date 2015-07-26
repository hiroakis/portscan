.PHONY: all clean deps build

all: clean build

deps:
	go get -d -v ./...
	go get github.com/mitchellh/gox

build:
	mkdir -p ./build
	gox -osarch="linux/amd64 darwin/amd64" -output ./build/portscan_{{.OS}}-{{.Arch}}

clean:
	rm -rf ./build
