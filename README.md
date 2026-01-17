# HTTP Server

This is a simple HTTP server built with Go. It is a basic server that can handle HTTP requests and responses.

## Features

- Handle HTTP requests and responses
- Recovery middleware:
    - It recovers from panics and returns a 500 error to the client
- Request ID middleware
    - It adds a unique ID to each request
- Logger middleware:
    - It logs the request and response
- Panic handler for testing recovery middleware:
    - It panics to test the recovery middleware
