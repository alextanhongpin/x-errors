# errors [![Go Reference](https://pkg.go.dev/badge/github.com/alextanhongpin/errors.svg)](https://pkg.go.dev/github.com/alextanhongpin/errors)

Better error handling for REST/GraphQL servers.


## User Story

As a developer, I want to avoid exposing the application errors to the user, so that users are not exposed to the internal errors.


## Context

Most REST applications (as well as GraphQL) requires mapping errors and returning them to the client. This errors are normally tied to the HTTP status code too, and hence requires mapping from domain errors to the corresponding status code 


However, the current `error` package has two main limitations

- it does not provide sufficient context for categorizing errors
- you cannot pass values through the errors 

WIP
