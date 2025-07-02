# ======================
# GoTasker Pro Makefile
# ======================

APP_NAME := gotasker-pro

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	pre-commit run golangci-lint --all-files

.PHONY: fmt
fmt:
	pre-commit run go-fmt --all-files

.PHONY: vet
vet:
	pre-commit run go-vet --all-files

.PHONY: check
check:
	pre-commit run --all-files

.PHONY: test
test:
	go test -v ./... -count=1

.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/api

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: install-hooks
install-hooks:
	pre-commit install
