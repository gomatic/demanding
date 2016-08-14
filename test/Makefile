export COMMIT_ID := $(shell git describe --tags --always --dirty 2>/dev/null)
export COMMIT_TIME := $(shell git show -s --format=%ct 2>/dev/null)

export GOARCH ?= amd64

test: test-$(shell uname | tr [A-Z] [a-z])

test-%: main.go Makefile
	go build -ldflags "-X main.VERSION=$(COMMIT_TIME)-$(COMMIT_ID)" -o test-$(GOOS)

image: GOOS := linux
image: test-$(GOOS) Dockerfile
	docker build --tag demanding/test:latest .

run: image
	docker run -d -p 3001:3001 --name demanding_test demanding/test

start:
	docker start demanding_test

stop:
	docker stop demanding_test

test: start
	curl localhost:3001/debug/vars
