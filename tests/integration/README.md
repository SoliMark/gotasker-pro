# Integration Tests

This directory contains integration tests for the GoTasker Pro application using Testcontainers.

## Test Structure

```
tests/integration/
├── suite.go          # Test suite setup and common utilities
├── user_test.go      # User-related integration tests
├── task_test.go      # Task-related integration tests
└── README.md         # This file
```

## Test Files

### `suite.go`
- Contains the `ContainerTestSuite` struct and setup methods
- Manages PostgreSQL and Redis containers using Testcontainers
- Provides common utilities like `createTestUserAndLogin`
- Handles database cleanup between tests

### `user_test.go`
- **TestUserRegistrationWithContainers**: Tests user registration
- **TestUserLoginWithContainers**: Tests user login
- **TestUserEndToEndFlowWithContainers**: Tests complete user flow
- **TestUserWithoutTokenWithContainers**: Tests unauthorized access
- **TestUserWithInvalidTokenWithContainers**: Tests invalid token handling

### `task_test.go`
- **TestTaskCRUDWithContainers**: Tests task CRUD operations
- **TestTaskAuthorizationWithContainers**: Tests task authorization
- **TestTaskValidationWithContainers**: Tests task validation
- **TestTaskCachingWithContainers**: Tests Redis caching functionality

## Running Tests

### Run all integration tests:
```bash
go test ./tests/integration/ -v
```

### Run specific test:
```bash
go test ./tests/integration/ -v -run TestUserRegistrationWithContainers
```

### Run all tests in a file:
```bash
go test ./tests/integration/ -v -run TestContainerTestSuite
```

## Test Environment

The integration tests use Testcontainers to spin up:
- **PostgreSQL 15**: Database for the application
- **Redis 7**: Cache layer
- **Testcontainers Ryuk**: Resource reaper for cleanup

Each test runs in isolation with a clean database state.

## Test Coverage

The integration tests cover:
- ✅ User authentication (register, login)
- ✅ Authorization (JWT tokens, protected routes)
- ✅ Task CRUD operations
- ✅ Task authorization (users can only access their own tasks)
- ✅ Data validation
- ✅ Redis caching functionality
- ✅ Error handling

## Dependencies

- `github.com/testcontainers/testcontainers-go`
- `github.com/testcontainers/testcontainers-go/modules/postgres`
- `github.com/testcontainers/testcontainers-go/modules/redis`
- `github.com/stretchr/testify/suite`
- `github.com/gin-gonic/gin`
