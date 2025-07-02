# GoTasker Pro

A practical backend side project built with **Go**, following real-world best practices for testing, linting, and CI/CD.

---

## 🚀 Usage

### 📦 Install dependencies
```go
go mod tidy
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
```go
make run
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
