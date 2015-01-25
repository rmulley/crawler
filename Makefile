GO ?= go
export GOPATH := $(CURDIR)/_vendor:$(GOPATH)

all: build

build:
	$(GO) fmt src/*
	$(GO) vet src/*
	$(GO) build -o ./bin/crawler src/main.go

linux:
	$(GO) fmt src/*
	$(GO) vet src/*	
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/crawler src/main.go

test:
	$(GO) test src/* -v -cover 

bench:
	$(GO) test src/* -v -bench=.

clean:
	rm -r _vendor/pkg/darwin_amd64/