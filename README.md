# Stress Test - CLI Load Testing Tool

A high-performance load testing tool written in Go that allows you to perform HTTP stress tests on web services with configurable concurrency and request counts.

## Features

- **CLI-based load testing** with simple parameter configuration
- **Concurrent request execution** with configurable worker count
- **Comprehensive reporting** including:
  - Total execution time
  - Request success rate
  - Status code distribution
  - Min/Max/Average response times
  - Requests per second throughput
- **Docker support** for easy deployment
- **Zero external dependencies** (uses Go standard library only)

## Installation

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized execution)

### Building from source

```bash
go build -o stress-test main.go
```

## Usage

### Command Line Parameters

- `--url`: URL of the service to test (required)
- `--requests`: Total number of requests to perform (required)
- `--concurrency`: Number of concurrent requests (required)

### Local Execution

```bash
# Basic example
./stress-test --url=http://example.com --requests=100 --concurrency=10

# Using Google as test target
./stress-test --url=https://www.google.com --requests=1000 --concurrency=50
```

### Docker Execution

```bash
# Build the Docker image
docker build -t markuscandido/stress-test-go .

# Run the test
docker run markuscandido/stress-test-go --url=http://example.com --requests=1000 --concurrency=10
```

## Example Output

```
Starting load test
URL: https://www.google.com
Total Requests: 100
Concurrency: 10

========== LOAD TEST REPORT ==========
Total Time:        3.245s
Total Requests:    100
Successful (200):  100
Min Duration:      245ms
Max Duration:      892ms
Avg Duration:      523ms
Requests/second:   30.82

Status Code Distribution:
  HTTP 200: 100
=======================================
```

## How It Works

1. **Concurrent Workers**: The tool spawns N workers (based on the `--concurrency` parameter)
2. **Request Distribution**: Each worker makes HTTP GET requests to the target URL until all requests are completed
3. **Thread-safe Tracking**: Uses atomic operations to track metrics without locks
4. **Report Generation**: After all requests complete, displays comprehensive statistics

## Report Metrics

- **Total Time**: Wall-clock time from start to finish
- **Total Requests**: Number of requests executed (should match --requests parameter)
- **Successful (200)**: Count of HTTP 200 responses
- **Min/Max/Avg Duration**: Response time statistics
- **Requests/second**: Throughput metric
- **Status Code Distribution**: Breakdown of all HTTP responses received

## Performance Notes

- Each worker uses an HTTP client with a 30-second timeout
- Concurrent execution ensures the test can reach high throughput
- Memory usage scales with the number of concurrent workers
- Network and target service performance will significantly affect results

## Error Handling

- Connection timeouts are tracked as "Error (timeout/failed)" in the report
- Failed requests don't stop the tool - it continues until all requests are attempted
- Network errors are counted separately from HTTP status codes

## License

MIT License
