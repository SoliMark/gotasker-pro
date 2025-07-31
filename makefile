# ================================
# GoTasker Pro Makefile
# All tasks run through pre-commit hooks
# ================================

APP_NAME := gotasker-pro

# ================================
# 1. Go Format, Lint, Test, Build
# ================================

.PHONY: tidy fmt vet lint check test build run clean mocks install-hooks

tidy:
	pre-commit run go-tidy --all-files

fmt:
	pre-commit run go-fmt --all-files

vet:
	pre-commit run go-vet --all-files

lint:
	pre-commit run golangci-lint --all-files

check:
	pre-commit run --all-files

test:
	pre-commit run go-test --all-files

build:
	pre-commit run go-build --all-files

run:
	go run ./cmd/api/main.go

clean:
	rm -rf bin/ coverage.out *.test
	@echo "Cleaned build artifacts."

# ================================
# 2. Mock Generation
# ================================

mocks:
	mockgen -source=internal/service/user_service.go \
		-destination=internal/service/mock_service/mock_user_service.go \
		-package=mock_service

	mockgen -source=internal/repository/user_repository.go \
		-destination=internal/repository/mock_repository/mock_user_repository.go \
		-package=mock_repository

	mockgen -source=internal/repository/task_repository.go \
  		-destination=internal/repository/mock_repository/mock_task_repository.go \
  		-package=mock_repository

# ================================
# 3. Pre-commit Hooks
# ================================

install-hooks:
	pre-commit install
	@echo "Pre-commit hooks installed."

# ================================
# 4. Docker Compose Controls
# ================================

.PHONY: up down downv restart logs

up:
	docker compose up -d
	@echo "Docker containers are up."

down:
	docker compose down
	@echo "Docker containers stopped (volumes kept)."

downv:
	docker compose down -v
	@echo "Docker containers and volumes removed."

restart: down up
	@echo "Docker containers restarted."

logs:
	docker compose logs -f
