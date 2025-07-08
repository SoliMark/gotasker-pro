# ======================
# GoTasker Pro Makefile
# 全部透過 pre-commit hooks 執行
# ======================

APP_NAME := gotasker-pro

.PHONY: tidy
tidy:
	pre-commit run go-tidy --all-files

.PHONY: fmt
fmt:
	pre-commit run go-fmt --all-files

.PHONY: vet
vet:
	pre-commit run go-vet --all-files

.PHONY: lint
lint:
	pre-commit run golangci-lint --all-files

.PHONY: check
check:
	pre-commit run --all-files

.PHONY: test
test:
	pre-commit run go-test --all-files

.PHONY: build
build:
	pre-commit run go-build --all-files

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: clean
clean:
	rm -rf bin/

.PHONY: mocks

mocks:
	mockgen -source=internal/service/user_service.go \
		-destination=internal/service/mock_service/mock_user_service.go \
		-package=mock_service

	mockgen -source=internal/repository/user_repository.go \
		-destination=internal/repository/mock_repository/mock_user_repository.go \
		-package=mock_repository


.PHONY: install-hooks
install-hooks:
	pre-commit install
