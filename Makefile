pwd = $(shell pwd)
GOPATH := $(pwd)/vendor:$(pwd)
export GOPATH

all:
	mkdir -p bin
	go build -o bin/gonc .

clean:
	rm -f bin/gonc

realclean: clean
	rm -rf vendor

