# Testing Framework for xcclookup Package

This document describes the testing framework implemented for the xcclookup package, including unit tests, integration tests, and testing utilities.

## Overview

The testing framework has been designed to provide comprehensive coverage of the xcclookup package functionality while maintaining good separation of concerns and testability.

## Architecture

### Refactoring Summary

The original code has been refactored to improve testability:

1. **Interface-based Design**: Created interfaces for external dependencies
2. **Service Layer**: Extracted business logic into service classes
3. **Dependency Injection**: Made components accept dependencies as parameters
4. **Mock Implementations**: Created mock objects for testing

### Key Components

#### Interfaces (`pkgs/types/interfaces.go`)

- `VersionServiceInterface`: For version operations
- `CacheManagerInterface`: For cache operations
- `ClaimsServiceInterface`: For claims operations
- `RPCClientInterface`: For RPC client operations
- `OutputFunction`: For output functions

#### Services (`pkgs/services/`)

- `VersionService`: Handles version-related operations
- `CacheService`: Handles cache-related operations
- `AssetsService`: Handles assets-related operations

#### Handlers (`pkgs/handlers/`)

- `VersionHandler`: HTTP handler for version requests
- `CacheHandler`: HTTP handler for cache requests
- `AssetsHandler`: HTTP handler for assets requests

#### Mocks (`pkgs/mocks/`)

- `MockVersionService`: Mock implementation of VersionServiceInterface
- `MockCacheManager`: Mock implementation of CacheManagerInterface
- `MockClaimsService`: Mock implementation of ClaimsServiceInterface

## Test Structure

### Unit Tests

#### Service Tests

- `services/version_service_test.go`: Tests for VersionService
- `services/cache_service_test.go`: Tests for CacheService
- `services/assets_service_test.go`: Tests for AssetsService

#### Handler Tests

- `handlers/version_handler_test.go`: Tests for VersionHandler
- `handlers/cache_handler_test.go`: Tests for CacheHandler
- `handlers/assets_handler_test.go`: Tests for AssetsHandler

#### Claims Tests

- `claims/liquid_test.go`: Tests for liquid claims functionality

### Integration Tests

- `integration/handlers_integration_test.go`: End-to-end tests for handlers

## Running Tests

### Using the Test Runner Script

```bash
./run_tests.sh
```

This script will:

1. Run all unit tests
2. Run integration tests
3. Generate coverage reports
4. Create an HTML coverage report

### Manual Test Execution

```bash
# Run all tests
go test -v ./pkgs/...

# Run specific package tests
go test -v ./pkgs/services/...
go test -v ./pkgs/handlers/...
go test -v ./pkgs/claims/...

# Run with coverage
go test -v -coverprofile=coverage.out ./pkgs/...
go tool cover -html=coverage.out -o coverage.html
```

## Test Patterns

### Mock Usage

```go
// Create mock
mockVersionService := &mocks.MockVersionService{
    GetVersionFunc: func() ([]byte, error) {
        return []byte(`{"version":"1.0.0"}`), nil
    },
}

// Use mock in service
service := services.NewVersionService(mockVersionService)
```

### Table-Driven Tests

```go
tests := []struct {
    name           string
    input          string
    expectedResult string
    expectedError  error
}{
    // Test cases...
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test implementation
    })
}
```

### HTTP Handler Testing

```go
req := httptest.NewRequest(http.MethodGet, "/version", nil)
w := httptest.NewRecorder()
handler.Handle(w, req)

assert.Equal(t, http.StatusOK, w.Code)
assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
```

## Coverage Goals

The testing framework aims for:

- **Unit Test Coverage**: >90% for business logic
- **Integration Test Coverage**: All major workflows
- **Error Handling**: All error paths tested
- **Edge Cases**: Boundary conditions and edge cases

## Best Practices

### Writing Tests

1. **Arrange-Act-Assert**: Structure tests with clear sections
2. **Descriptive Names**: Use descriptive test and variable names
3. **Single Responsibility**: Each test should test one thing
4. **Mock External Dependencies**: Don't test external services in unit tests
5. **Test Error Conditions**: Always test error paths

### Test Data

1. **Use Constants**: Define test data as constants
2. **Minimal Test Data**: Use only necessary data for each test
3. **Realistic Data**: Use realistic but minimal test data

### Assertions

1. **Specific Assertions**: Use specific assertions rather than generic ones
2. **Meaningful Messages**: Provide meaningful assertion messages
3. **Multiple Assertions**: Test multiple aspects when appropriate

## Maintenance

### Adding New Tests

1. Follow the existing patterns
2. Use the appropriate mock interfaces
3. Add integration tests for new workflows
4. Update this documentation

### Updating Tests

1. Keep tests in sync with code changes
2. Update mocks when interfaces change
3. Maintain test coverage goals
4. Review and refactor tests regularly

## Troubleshooting

### Common Issues

1. **Import Errors**: Ensure all dependencies are properly imported
2. **Mock Issues**: Verify mock implementations match interfaces
3. **Test Failures**: Check that test data matches expected results
4. **Coverage Issues**: Add tests for uncovered code paths

### Debugging Tests

1. Use `go test -v` for verbose output
2. Add debug prints in test functions
3. Use `t.Log()` for test logging
4. Check test coverage reports for gaps

## Future Enhancements

1. **Performance Tests**: Add benchmarks for critical paths
2. **Property-Based Testing**: Use property-based testing for complex logic
3. **Contract Testing**: Add contract tests for external dependencies
4. **Test Data Factories**: Create factories for complex test data
5. **Parallel Testing**: Enable parallel test execution where appropriate
