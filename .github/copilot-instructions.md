# Stress Test - Load Testing CLI

This is a Go CLI application for performing HTTP load tests with configurable concurrency and request counts.

## Project Overview

- **Language**: Go 1.21+
- **Type**: CLI Application
- **Build**: Standard Go build system
- **Docker**: Included with multi-stage build
- **Dependencies**: None (standard library only)

## Building & Running

### Build locally:
```bash
go build -o stress-test main.go
./stress-test --url=http://example.com --requests=100 --concurrency=10
```

### Docker:
```bash
docker build -t stress-test .
docker run stress-test --url=http://example.com --requests=1000 --concurrency=10
```

## Key Features

- Concurrent HTTP request execution
- Real-time progress tracking
- Comprehensive performance report
- Support for custom concurrency levels
- Error handling for network issues

## File Structure

- `main.go`: Core application code with load test logic
- `go.mod`: Go module definition
- `Dockerfile`: Multi-stage Docker build
- `README.md`: User documentation

## Next Steps

Test the application locally or build and push the Docker image to your registry.
