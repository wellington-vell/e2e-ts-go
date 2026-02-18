# e2e-ts-go

An end-to-end type-safe full-stack project inspired by [Better T Stack](https://www.better-t-stack.dev/), featuring a **Go backend** and **TypeScript frontend** with full type safety across the stack.

## Overview

This project demonstrates how to achieve end-to-end type safety between a Go backend and TypeScript frontend, bringing the developer experience of modern TypeScript full-stack frameworks to a Go-based architecture.

## Implementations

Different type safety implementations are in branches. Each branch explores a different approach or pattern for achieving end-to-end type safety:

- **open-api** - Custom implementation using [oRPC](https://orpc.unnoq.com/) approach with OpenAPI
- **connect/grpc** - Implementation using gRPC with Protocol Buffers

Switch branches to explore different implementations:

```bash
git checkout open-api   # Custom oRPC implementation with OpenAPI
git checkout connect/grpc   # gRPC implementation
```

## Opinion

### gRPC Experience

After working with the gRPC implementation, I found it to be unnecessarily complex for most web applications:

- **Steep learning curve**: Had to learn Protocol Buffers, gRPC concepts, and code generation workflows - a whole new paradigm just for type safety
- **Tooling overhead**: Extra protoc compiler, language-specific plugins, and generated code management
- **Development friction**: Hot reload becomes painful with generated files, and debugging is harder with binary payloads
- **Limited browser support**: gRPC-Web adds another layer of complexity with its own proxy requirements
- **Overkill for simple APIs**: Most CRUD operations don't need streaming or the performance benefits gRPC provides

### Custom oRPC implementation with OpenAPI Experience

The oRPC approach feels more natural:

- **Familiar HTTP/REST patterns**: No paradigm shift required
- **Better tooling ecosystem**: Standard OpenAPI tools work out of the box
- **Easier debugging**: Human-readable JSON over the wire

**Verdict**: For most full-stack TypeScript + Go projects, with OpenAPI provides 90% of the type safety benefits with 10% of the complexity. gRPC might make sense for high-performance microservices, but it's overkill for typical web applications.

## Goals

- ✅ Go backend with HTTP server
- ✅ TypeScript frontend with modern React stack
- ✅ Monorepo setup with Turbo
- 🚧 End-to-end type safety between Go and TypeScript
- 🚧 Type generation from Go to TypeScript
- 🚧 RPC-like communication with type safety

## Inspiration

This project is inspired by [Better T Stack](https://www.better-t-stack.dev/), which provides a modern CLI for scaffolding end-to-end type-safe TypeScript projects. This implementation explores achieving similar type safety guarantees with a Go backend instead of a TypeScript/Node.js backend.
