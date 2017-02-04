# Package Dependency Management Challenge : 0.1.2

## Features

1. All Tests Passing with regard to test harness.
2. Modularized solution with 4 packages : data/err/input/operation
3. Testing files in place beside main files in each package.
4. Concurrency of 100 clients reached, locally.  Maximum is currently 32 for Docker image.
5. Optional throttling of requests for rate limiting.
6. Dockerized service, and docker-compose.yml included up for easy local deployment.


## Building and Running

### Local Host

<pre><code>
go build
go run main.go [-throttle 1000]
</code></pre>

### Docker Compose

<pre><code>
docker-compose build
docker-compose up -d
docker-compose logs -f
</code></pre>

### Just Docker

<pre><code>
docker build -t pkgindexer .
docker run -p 8080:8080 -d pkgindexer
</code></pre>

## Future Plans

- Input validation improvement using regex and sanitization.
- Business logic algorithm optimized, DRY principle, made readable.
- In memory data store and optimized locking - do we really need to lock the entire in memory cache?
- Error handling using custom error types.
- Uniform logging system.
- Documentation generated and included.
- Automated integration testing in addition to unit tests that have been added, modeled after test harness.
- Externalization of environment configuration.
- Basic performance tests.
- Basic linting for code consistency.
  
## Versions
- 0.1.0 - Initial messy solution to pass test harness with low concurrency.
- 0.1.1 - Modularized solution split out between packages.
- 0.1.2 - Concurrency goals met for local run, optional throttling implemented.

## Performance Notes

### Request Throttling Comparisons

* Utilizing the Docker image
* Concurrency = 32
* 3 runs per configuration.

| Rate  | Test Harness Run Time  |
|--:|---|
| Unthrottled  | 4.7 - 4.9s |
| 1000 requests/client/second | 4.7-5 seconds  |
| 500 requests/client/second  | 5.5-5.8 seconds  |
| 200 requests/client/second  | 11.2-12.4 seconds  |
| 100 requests/client/second  | 24.1-24.9 seconds  |

