.PHONY: run dev binary setup glide test update
SHELL := /bin/bash

all: run

run: binary
	./frontend

dev:
	rm -f gin-bin || true
	gin --all run main.go

binary:
	GOOS=linux go build -i -o frontend

setup:
	go get -v -u github.com/codegangsta/gin
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)