# GoRPC

GoRPC is a minimal RPC framework for Go, inspired by [oRPC](https://orpc.dev/docs/getting-started). It provides a simple, type-safe way to build RPC services using Go's standard library, with a custom radix tree router for efficient path matching.

## Overview

GoRPC combines RPC (Remote Procedure Call) patterns with HTTP, allowing you to define and call remote procedures through a clean, type-safe API. The framework is designed to be minimal and lightweight.

## Features

- **Custom radix tree router** - Efficient O(n) path matching with support for multiple path parameters at any position
- **Path parameter extraction** - Automatic extraction of path parameters (e.g., `:id`, `:userId`) into `Context.Params`
- **Flexible routing** - Supports routes like `/users/:userId/posts/:postId` with parameters at any position
- **Context-based handlers** with access to request, response, and path parameters
- **JSON serialization** for input/output
- **Error handling** with HTTP status codes (400, 404, 500, etc.)
- **Middleware system** for request/response processing
- **Procedure builder** with fluent `OS()` builder for type-safe procedure creation
- **Metadata support** with Meta, Tags, and Route definitions for procedures
- **OpenAPI integration** with automatic OpenAPI 3.0 spec generation
- **Plugin system** with Scalar and Swagger UI plugins
- **CORS support** for cross-origin requests
- **Minimal API surface** for easy adoption

## Architecture

GoRPC uses a custom radix tree (prefix tree) router implementation for route registration and matching. This provides:

- Efficient path matching in O(n) time where n is the number of path segments
- Support for multiple path parameters at any position (e.g., `/api/users/:userId/posts/:postId`)
- Automatic parameter extraction and storage in `Context.Params`
- Proper HTTP status code handling (404 for not found, 405 for method not allowed)

Each procedure is registered as an HTTP endpoint with support for RESTful path patterns.

### Core Components

- **GORPC** - Main application instance
- **Router** - Map of procedure names to procedures
- **Procedure** - Individual RPC procedure definition
- **ProcedureBuilder** - Fluent builder for creating procedures with `OS[TInput, TOutput]()`
- **Context** - Request context with params, headers, and body
- **Middleware** - Request/response processing pipeline
- **Plugin** - Extensible plugin system (OpenAPI, Scalar, Swagger)

### API Methods

- **`New()`** - Create a new GORPC instance
- **`Prefix()`** - Add a prefix to all routes
- **`Router()`** - Register HTTP handlers
- **`Middleware()`** - Add global middleware
- **`Plugin()`** - Register plugins (OpenAPI, Scalar, Swagger)
- **`EnableCORS()`** - Enable CORS support
- **`ListenAndServe()`** - Start the HTTP server

### Procedure Builder

- **`OS[TInput, TOutput]()`** - Create a new procedure builder
- **`Tag()`** - Add tags for OpenAPI documentation
- **`Meta()`** - Add metadata (summary, description)
- **`Route()`** - Define HTTP method and path
- **`Errors()`** - Define error status codes
- **`Handler()`** - Register the handler function
- **`Build()`** - Build the procedure

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/wellington-vell/gorpc"
)

type Input struct {
	Name string `json:"name"`
}

type Output struct {
	Message string `json:"message"`
}

var router = gorpc.Router{
	"greet": gorpc.OS[Input, Output]().
		Meta(gorpc.Meta{
			Summary:     "Greet endpoint",
			Description: "Returns a greeting message",
		}).
		Route(gorpc.Route{
			Method: "POST",
			Path:   "/greet/:name",
		}).
		Handler(func(ctx *gorpc.Context, input Input) (Output, error) {
			return Output{Message: "Hello, " + input.Name}, nil
		}).
		Build(),
}

func main() {
	app := gorpc.New().
		EnableCORS(gorpc.CORSConfig{
			Origins:     []string{"*"},
			Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			Headers:     []string{"Content-Type", "Authorization"},
			Credentials: true,
		}).
		Prefix("/api").
		Router(router).
		Plugin(gorpc.NewScalarPlugin()).
		Plugin(gorpc.NewSwaggerPlugin())

	addr := ":8080"
	fmt.Printf("Server starting on %s\n", addr)

	if err := app.ListenAndServe(addr); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}
```

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

## OpenAPI Integration

GoRPC automatically generates OpenAPI 3.0 specifications for all registered procedures. Access the spec at `/openapi.json`.

Use the Scalar plugin for interactive API documentation:

```go
app.Plugin(gorpc.NewScalarPlugin())
```

## Standard Library Only

This package adheres to a strict "standard library only" policy. All functionality is implemented using packages from the Go standard library (`net/http`, `encoding/json`, `context`, `reflect`, etc.), ensuring no external dependencies and maximum compatibility. The radix tree router is implemented from scratch using only standard library packages.

## Acknowledgments

GoRPC is inspired by:

- **[go-chi/chi](https://github.com/go-chi/chi)** - Lightweight, idiomatic router for building Go HTTP services. GoRPC's radix tree router implementation draws inspiration from chi's approach to path matching and middleware patterns.

- **[oRPC](https://orpc.dev/)** - Type-safe API framework for TypeScript. GoRPC follows oRPC's philosophy of providing a simple, type-safe way to build RPC services with a focus on developer experience.
