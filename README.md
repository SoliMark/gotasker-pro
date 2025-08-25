# GoTasker Pro

A practical backend side project built with **Go**, following real-world best practices for testing, linting, and CI/CD.

## 🎯 Project Overview

GoTasker Pro 是一個任務管理系統，採用 **Layered Architecture** 設計，支援用戶認證、任務 CRUD 操作和 Redis 快取功能。

### ✨ Features
- 🔐 JWT 用戶認證系統
- 📝 完整的任務 CRUD 操作
- ⚡ Redis 快取支援（Cache-Aside Pattern）
- 🧪 完整的單元測試覆蓋
- 🐳 Docker 容器化支援
- 🔄 CI/CD 自動化流程

### 🏗️ Architecture
```
┌─────────────────┐
│   HTTP Layer    │  ← Gin Router + Handlers
├─────────────────┤
│  Service Layer  │  ← Business Logic + Cache
├─────────────────┤
│ Repository Layer│  ← Data Access (GORM)
├─────────────────┤
│   Model Layer   │  ← Domain Models
└─────────────────┘
```

---

## 🚀 Quick Start

### 📋 Prerequisites
- Go 1.23+
- PostgreSQL 或 SQLite
- Redis (可選，用於快取)

### 📦 Install dependencies
```bash
go mod tidy
```

### ⚙️ Configuration
複製環境變數範例文件：
```bash
cp env.example .env
```

編輯 `.env` 文件，設定必要的環境變數：
```bash
# 必要配置
DB_URL=postgres://username:password@localhost:5432/gotasker_pro?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-here

# 可選配置（Redis 快取）
REDIS_ADDR=localhost:6379
CACHE_TTL_TASKS=60s
```

### 🐳 Using Docker Compose
啟動 PostgreSQL 和 Redis：
```bash
docker-compose up -d
```

### ✅ Run lint (uses pre-commit golangci-lint)

```go
make lint
```

### 🧹 Run format (gofmt)
```go
make fmt
```

### 🔍 Run vet
```go
make vet
```

### ✅ Run all pre-commit hooks

```go
make check
```
### 🧪 Run unit tests
```go
make test
```

### ⚙️ Build binary
```go
make build
```

### 🏃 Run application
```bash
make run
```

### 📚 API Documentation
詳細的 API 文檔請參考：[docs/API.md](docs/API.md)

### 🧪 Testing
運行所有測試：
```bash
make test
```

運行特定測試：
```bash
go test ./internal/service -v
```

### ⚡️ Pre-commit hooks
This project uses pre-commit to automate linting and formatting before each commit.

### To install hooks locally:

```go
make install-hooks
```

### To run all hooks manually:
```go
make check
```

### 📋 Notes
Linting is handled by golangci-lint via pre-commit.
