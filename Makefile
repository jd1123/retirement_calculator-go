export GOPATH:=$(shell pwd)

default: all

deps:
	go get -d -v retirement_calculator-go

build:
	go build retirement_calculator-go

test:
	cd src/retirement_calculator-go
	go test

clean-build:
	rm -rf retirement_calculator-go

# There is a better way to do this
clean-deps:
	#cd src/retirement_calculator-go
	#go clean -i -r -n
	rm -rf src/bitbucket.org
	rm -rf src/github.com
	rm -rf src/golang.org
	rm -rf pkg

clean-all:
	make clean-build
	make clean-deps

all:
	make clean-build
	make deps
	make build
