# Package Dependency Management Challenge : 0.1.1

## Features

1. All Tests Passing with regard to test harness with low concurrency.
2. Modularized solution with 4 packages : data/err/input/operation
3. Testing files in place beside main files in each package.
4. Concurrency of 130 reached

## Future Plans

### Required
- Connection management : Rate limiting, identifying spammy clients.
- Input validation improvement using regex and sanitization.
- Business logic algorithm optimized, DRY principle, made readable.
- In memory data store and optimized locking - do we really need to lock the entire in memory cache?
- Error handling using custom error types.
- Uniform logging system.
- Documentation generated and included.
- Automated integration testing in addition to unit tests that have been added, modeled after test harness.
- Externalization of environment configuration.
- Docker-ization of service.
- Basic performance tests.
- Basic linting for code consistency.
  
## Versions
- 0.1.0 - Initial messy solution to pass test harness with low concurrency.
- 0.1.1 - Modularized solution split out between packages.
