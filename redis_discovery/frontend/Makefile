.PHONY: run binary setup glide test update
SHELL := /bin/bash

all: run

run: binary
	./frontend

binary:
	GOOS=linux go build -i -o frontend

setup:
	go get -v -u github.com/Masterminds/glide

glide:
	glide install --force

test:
	GOARCH=amd64 GOOS=linux go test $$(go list ./... | grep -v /vendor/)