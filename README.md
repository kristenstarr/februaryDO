# Package Dependency Management Challenge : 0.1.4

## Features

1. All Tests Passing with regard to test harness.
2. Modularized solution with 4 packages : data/err/input/operation
3. Testing files in place beside main files in each package, complete unit testing.
4. Concurrency of 100 clients reached, locally.  Maximum is currently 32 for Docker image.
5. Optional throttling of requests for rate limiting.
6. Input validation utilizing Regex.
7. Dockerized service, and docker-compose.yml included up for easy local deployment.


## Building and Running

### Local Host

<pre><code>go build
go run main.go [-throttle 1000]
</code></pre>

### Docker Compose

<pre><code>docker-compose build
docker-compose up -d
docker-compose logs -f
</code></pre>

### Just Docker

<pre><code>docker build -t pkgindexer .
docker run -p 8080:8080 -d pkgindexer
</code></pre>

## Future Plans

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
- 0.1.3 - Input validation and tests.
- 0.1.4 - Cleaning up and testing of business logic.

## Performance Notes

### Request Throttling Comparisons

* Utilizing the Docker image
* Concurrency = 32
* 3 runs per configuration.

| Rate  | Test Harness Run Time  |
|---|---|
| Unthrottled  | 4.7 - 4.9s |
| 1000 requests/client/second | 4.7-5 seconds  |
| 500 requests/client/second  | 5.5-5.8 seconds  |
| 200 requests/client/second  | 11.2-12.4 seconds  |
| 100 requests/client/second  | 24.1-24.9 seconds  |

## Future Considerations

### Locking

Currently, the entire in-memory cache representing the Index is locked for each request that comes in.
This is done so that we do not have concurrency issues with reading/writing the same keys in the map,
since many libraries and their dependencies are interconnected.
Consideration was put into whether a separate locking system should be built - locking only on a library
and its dependencies, but there are a few reasons why this was decided against :
1. In order to obtain a lock specific to a set of libraries, we would have to check if a lock already exists
for this particular set of libraries, if not create one - and this process would have to be synchronized
thus pushing the locking step only one step up. If we were using a persistent datastore (sql/nosql)
and could use in-memory for these locks, then we could indeed speed up the process by a universal
lock on 'lock creation and obtaining' and only do partial locks on the slower data calls.
2. For such operations as Remove, we are unsure until we query for the Library's parents which other
libraries we'd need to lock on.  We don't want to be removing a library from another library's parent
list while that library itself is being removed.
3. The chosen implementation is the most 'transaction safe', and although it perhaps causes some unnecessary
waiting (in the case where libraries being requested are not connected at all), it is a simple and clean
implementation that seems to fit well for the widely-interconnected domain that we are working with.
