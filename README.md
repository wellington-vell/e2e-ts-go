# e2e-ts-go

An end-to-end type-safe full-stack project inspired by [Better T Stack](https://www.better-t-stack.dev/), featuring a **Go backend** and **TypeScript frontend** with full type safety across the stack.

## Overview

This project demonstrates how to achieve end-to-end type safety between a Go backend and TypeScript frontend, bringing the developer experience of modern TypeScript full-stack frameworks to a Go-based architecture.

## Architecture

- **Backend**: Go server with HTTP handlers
- **Frontend**: TypeScript/React with Vite, TanStack Router, and TanStack Query
- **Monorepo**: Managed with Turbo for efficient builds and development
- **Type Safety**: End-to-end type safety between Go and TypeScript (work in progress)

## Tech Stack

### Backend
- **Go** - High-performance backend server
- **godotenv** - Environment variable management

### Frontend
- **React 19** - UI framework
- **TypeScript** - Type-safe frontend
- **Vite** - Fast build tool and dev server
- **TanStack Router** - Type-safe routing
- **TanStack Query** - Data fetching and state management
- **Tailwind CSS** - Utility-first CSS framework
- **Zod** - Runtime type validation

### Tooling
- **Turbo** - Monorepo build system
- **Bun** - Package manager and runtime
- **oxlint** - Fast TypeScript linter
- **oxfmt** - Code formatter

## Project Structure

```
e2e-ts-go/
├── apps/
│   ├── server/          # Go backend server
│   │   ├── cmd/         # Application entry point
│   │   └── internal/    # Internal packages
│   └── web/             # TypeScript frontend
│       └── src/
│           ├── components/
│           ├── lib/
│           └── routes/
├── package.json         # Root package.json for monorepo
└── turbo.json          # Turbo configuration
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.25.5 or later
- [Bun](https://bun.sh/) 1.3.5 or later

### Installation

```bash
# Install dependencies
bun install
```

### Development

```bash
# Start both server and web app in development mode
bun run dev
```

This will start:
- Go server (default port from `SERVER_PORT` env variable)
- Vite dev server for the frontend (default port from `VITE_WEB_PORT` env variable)

### Environment Variables

Create a `.env` file in the root directory:

```env
SERVER_PORT=3001
VITE_WEB_PORT=3000
```

## Goals

- ✅ Go backend with HTTP server
- ✅ TypeScript frontend with modern React stack
- ✅ Monorepo setup with Turbo
- 🚧 End-to-end type safety between Go and TypeScript
- 🚧 Type generation from Go to TypeScript
- 🚧 RPC-like communication with type safety

## Inspiration

This project is inspired by [Better T Stack](https://www.better-t-stack.dev/), which provides a modern CLI for scaffolding end-to-end type-safe TypeScript projects. This implementation explores achieving similar type safety guarantees with a Go backend instead of a TypeScript/Node.js backend.

## License

See [LICENSE](./LICENSE) file for details.
