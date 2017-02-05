#  Package Dependency Management Challenge : 0.1.7

##  Features

1. All Tests Passing with regard to test harness.
2. Modularized solution with 4 packages : data/err/input/operation
3. Testing files in place beside main files in each package, complete unit testing.
4. Concurrency of 100 clients reached, locally.  Maximum is currently 32 for Docker image.
5. Optional throttling of requests for rate limiting.
6. Input validation utilizing Regex.
7. Dockerized service, and docker-compose.yml included up for easy local deployment.
8. Simple custom logger that includes level (trace/info/error/debug)
9. Documentation on all modules, to be generated using godoc.
10. Integration and benchmark tests that utilize test client from provided test harness.


##  Building and Running

###  Local Host

<pre><code>go build
go run main.go
</code></pre>


Running locally is simplest, but our environment is not necessarily reproducible or consistent.

###  Docker Compose

<pre><code>docker-compose build
docker-compose up -d
docker-compose logs -f
</code></pre>


###  Just Docker

<pre><code>docker build -t pkgindexer .
docker run -p 8080:8080 -d pkgindexer
</code></pre>


##  Configuration


###  Throttling
By setting the 'throttle' value, we can limit each client in its ability to send messages to our
service at a capped rate.  Rate is given as an integer in messages per second per client.

<pre><code>go run main.go -throttle 1000</code></pre>

NOTE: Throttling is better observed by using the docker setups, as the environment is cleaner and
more reproducable than local environments.  See below Request Throttling Comparisons for some numbers
that make sense.


###  Logging
Log level can be set by using the logLevel parameter, which defaults to INFO.

<pre><code>go run main.go -logLevel TRACE</code></pre>

NOTE: Too much intensive logging, TRACE/DEBUG, will likely cause undesirable performance under load.


##  Testing
Automated tests include unit, functional, and benchmark tests. Unit and functional tests can be
run via the following from the main directory.

<pre><code>go test ./...
</code></pre>

NOTE : Integration tests are meant to be end-to-end tests of expected behavior.
Integration test module requires that a server be running on port 8080
that can be connected to by clients for validation.  If no server is running, then integration
tests will be skipped and only unit tests will be run.

Output should resemble the following

<pre>
?   	github.com/kristenfelch/pkgindexer	[no test files]
ok  	github.com/kristenfelch/pkgindexer/data	0.009s
ok  	github.com/kristenfelch/pkgindexer/err	0.010s
ok  	github.com/kristenfelch/pkgindexer/input	0.047s
ok  	github.com/kristenfelch/pkgindexer/integration	0.050s
?   	github.com/kristenfelch/pkgindexer/logging	[no test files]
ok  	github.com/kristenfelch/pkgindexer/operation	0.010s
</pre>

In order to include benchmark tests, use the following command instead.

<pre><code>go test ./... -bench=.
</code></pre>

Similar to functional tests, server must be running on port 8080.

##  Documentation
To generate and view logs, use the godoc module as such, with port of your choice:

<pre><code>godoc -http=:6060</code></pre>

Then you will be able to navigate to the following link to view documentation.

<pre><code>http://localhost:6060/pkg/github.com/kristenfelch/pkgindexer/</code></pre>
  

##  Versions
- 0.1.0 - Initial messy solution to pass test harness with low concurrency.
- 0.1.1 - Modularized solution split out between packages.
- 0.1.2 - Concurrency goals met for local run, optional throttling implemented.
- 0.1.3 - Input validation and tests.
- 0.1.4 - Cleaning up and testing of business logic.
- 0.1.5 - Custom error message type.
- 0.1.6 - Simple custom logging interface, and documentation added to work with godoc.
- 0.1.7 - Functional and benchmark tests


##  Performance Notes


###  A Few Benchmarks

One message of each type was chosen as a sampling of benchmark tests.  It is observable
that error messages are returned faster than others as they fail validation and require no
further processing.

| Name  | Number of Runs  | Average Run Time  |
|---|---|---|
| BenchmarkErrorMessages-8  | 30000  | 38065 ns/op  |
| BenchmarkIndexMessage-8  | 20000  | 63296 ns/op  |
| BenchmarkQueryMessage-8  | 20000  | 59702 ns/op  |
| BenchmarkRemoveNonIndexedMessage-8  | 20000  | 59586 ns/op  |


###  Request Throttling Comparisons

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


##  Future Considerations


###  Locking

Currently, the entire in-memory cache representing the Index is locked for each request that comes in.
This is done so that we do not have concurrency issues with reading/writing the same keys in the map,
since many packages and their dependencies are interconnected.
Consideration was put into whether a separate locking system should be built - locking only on a package
and its dependencies, but there are a few reasons why this was decided against.

1. In order to obtain a lock specific to a set of packages, we would have to check if a lock already exists
for this particular set of packages, if not create one - and this process would have to be synchronized
thus pushing the locking step only one step up. If we were using a persistent datastore (sql/nosql)
and could use in-memory for these locks, then we could indeed speed up the process by a universal
lock on 'lock creation and obtaining' and only do partial locks on the slower data calls.

2. For such operations as Remove, we are unsure until we query for the Package's parents which other
packages we'd need to lock on.  We don't want to be removing a package from another package's parent
list while that package itself is being removed.

3. The chosen implementation is the most 'transaction safe', and although it perhaps causes some unnecessary
waiting (in the case where packages being requested are not connected at all), it is a simple and clean
implementation that seems to fit well for the widely-interconnected domain that we are working with.


###  Concurrency: Docker versus Local
Running locally, it is easy to achieve concurrency at 100 clients.  When a docker image is spun up,
either using docker-compose or not, for some reason the max concurrency that can be used is 32. 
Open file limits and docker parameters have been investigated in attempts to solve this discrepancy,
but no solution has yet been reached.
