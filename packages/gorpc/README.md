# GoRPC

GoRPC is a minimal RPC framework for Go, inspired by [oRPC](https://orpc.dev/docs/getting-started). It provides a simple, type-safe way to build RPC services using Go's standard library, with a custom radix tree router for efficient path matching.

## Overview

GoRPC combines RPC (Remote Procedure Call) patterns with HTTP, allowing you to define and call remote procedures through a clean, type-safe API. The framework is designed to be minimal and lightweight, using only Go's standard library packages.

## Current Features

- **Custom radix tree router** - Efficient O(n) path matching with support for multiple path parameters at any position
- **Path parameter extraction** - Automatic extraction of path parameters (e.g., `:id`, `:userId`) into `Context.Params`
- **Flexible routing** - Supports routes like `/users/:userId/posts/:postId` with parameters at any position
- **Context-based handlers** with access to request, response, and path parameters
- **JSON serialization** for input/output
- **Error handling** with HTTP status codes (404, 405, etc.)
- **Minimal API surface** for easy adoption

## Architecture

GoRPC uses a custom radix tree (prefix tree) router implementation for route registration and matching. This provides:
- Efficient path matching in O(n) time where n is the number of path segments
- Support for multiple path parameters at any position (e.g., `/api/users/:userId/posts/:postId`)
- Automatic parameter extraction and storage in `Context.Params`
- Proper HTTP status code handling (404 for not found, 405 for method not allowed)

Each procedure is registered as an HTTP endpoint with support for RESTful path patterns.

### Core Features

- **Middleware system** - No middleware chain or execution pipeline
- **Procedure builder** - No `OS()` builder function for fluent procedure creation
- **Router builder** - No `NewRouterBuilder()` for router-level configuration
- **Input/Output type validation** - No automatic type checking or conversion
- **Custom validation** - No support for custom validation functions
- **Metadata support** - No Meta, Tags, or Route definitions for procedures
- **OpenAPI integration** - No OpenAPI spec generation or integration
- **Error maps** - No type-safe error handling with error code maps

### API Methods

- **`Use()` method** - Alias for `Router()` not available
- **`Serve()` method** - Method to return `http.Handler` not implemented
- **`GetRouters()` method** - Access to registered routers not exposed
- **`GetBasePath()` method** - Getter for base path not available
- **`GetRouterMapForOpenAPI()` method** - OpenAPI router map conversion not implemented

### Procedure Features

- **Input type conversion** - No automatic conversion from JSON to struct types
- **Output validation** - No validation of output types
- **Status code handling** - No automatic status code selection (e.g., 201 for create operations)
- **Context type support** - No typed context support beyond basic request/response
- **ProtectedProcedure** - No distinction between public and protected procedures

### Advanced Features

- **File upload/download** - Not supported
- **Server-Sent Events (SSE)** - Not supported
- **Request/Response interceptors** - Not supported
- **Plugin system** - Not supported
- **Client generation** - Not supported

## Design Philosophy

GoRPC follows a minimal design philosophy, focusing on core RPC functionality without external dependencies. The package uses only Go's standard library, making it lightweight and easy to integrate into existing projects.

## Routing

The custom radix tree router supports:
- **Exact path matching**: `/api/todos` matches exactly
- **Single parameter**: `/api/todos/:id` matches `/api/todos/123`
- **Multiple parameters**: `/api/users/:userId/posts/:postId` matches `/api/users/123/posts/456`
- **Parameters at any position**: `/api/todos/:id/comments/:commentId` works correctly

Path parameters are automatically extracted and available in `Context.Params` as a `map[string]string`. For example, a route `/api/users/:userId/posts/:postId` will populate `ctx.Params` with:
```go
ctx.Params["userId"] = "123"
ctx.Params["postId"] = "456"
```

## Standard Library Only

This package adheres to a strict "standard library only" policy. All functionality is implemented using packages from the Go standard library (`net/http`, `encoding/json`, `context`, `reflect`, etc.), ensuring no external dependencies and maximum compatibility. The radix tree router is implemented from scratch using only standard library packages.
