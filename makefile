# ======================
# GoTasker Pro Makefile
# ======================

APP_NAME := gotasker-pro

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./... -count=1

.PHONY: bulid
build :
	go build -o bin/$(APP_NAME) ./cmd/api

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: clean
clean:
	rm -rf bin/