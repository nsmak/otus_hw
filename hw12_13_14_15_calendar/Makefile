BIN_CAL := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./internal/server/grpcsrv/server.go

build-calendar:
	go build -v -o $(BIN_CAL) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

run: build-calendar
	$(BIN_CAL) -config ./configs/calendar.json

run: build-calendar
	$(BIN_SCHEDULER) -config ./configs/scheduler.json

run: build-sender
	$(BIN_SENDER) -config ./configs/sender.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN_CAL) version

up:
	sh ./deployments/deploy.sh up

down:
	sh ./deployments/deploy.sh down

test:
	go test -race ./internal/... ./pkg/...

integration-tests:
	sh ./deployments/deploy.sh tests

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
