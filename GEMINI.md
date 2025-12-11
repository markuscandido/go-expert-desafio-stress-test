# GEMINI.md

## Project Overview

This project is a high-performance, CLI-based load testing tool written in Go. It's designed to perform HTTP stress tests on web services with configurable concurrency and request counts. The tool is self-contained, using only the Go standard library, and can be run as a standalone executable or as a Docker container.

**Key Technologies:**

*   **Language:** Go (version 1.21 or higher)
*   **Build System:** Standard Go build tools
*   **Containerization:** Docker (multi-stage build)
*   **Dependencies:** None (Go standard library only)

**Architecture:**

The application spawns a number of concurrent workers (goroutines) to send HTTP GET requests to a target URL. It uses atomic operations for thread-safe metric tracking to ensure efficient and accurate reporting without the need for locks. After all requests are completed, it generates a comprehensive report on the performance, including total time, success rate, status code distribution, and response time statistics.

## Building and Running

### Building from Source

To build the executable from the source code, run the following command:

```bash
go build -o stress-test main.go
```

### Local Execution

To run the load test locally, use the compiled executable with the following parameters:

*   `--url`: The URL of the service to test (required).
*   `--requests`: The total number of requests to perform (required).
*   `--concurrency`: The number of concurrent requests to use (required).

**Example:**

```bash
./stress-test --url=http://example.com --requests=100 --concurrency=10
```

### Docker Execution

The project includes a multi-stage `Dockerfile` for containerized execution.

**1. Build the Docker Image:**

```bash
docker build -t stress-test .
```

**2. Run the Test in a Container:**

```bash
docker run stress-test --url=http://example.com --requests=1000 --concurrency=10
```

## Development Conventions

*   **Coding Style:** The code follows standard Go formatting and conventions.
*   **Dependencies:** The project intentionally has no external dependencies and relies solely on the Go standard library.
*   **Testing:** While there are no formal test files (`*_test.go`) in the project, testing is done via ad-hoc local or Docker-based execution against a target URL.
*   **Error Handling:** The tool is designed to be resilient. Network errors or non-200 responses are tracked and reported but do not halt the test.
