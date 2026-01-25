# e2e-ts-go

An end-to-end type-safe full-stack project inspired by [Better T Stack](https://www.better-t-stack.dev/), featuring a **Go backend** and **TypeScript frontend** with full type safety across the stack.

## Overview

This project demonstrates how to achieve end-to-end type safety between a Go backend and TypeScript frontend, bringing the developer experience of modern TypeScript full-stack frameworks to a Go-based architecture.

## Architecture

- **Backend**: Go server with HTTP handlers
- **Frontend**: TypeScript/React with Vite, TanStack Router, and TanStack Query
- **Monorepo**: Managed with Moonrepo for efficient builds and development
- **Type Safety**: End-to-end type safety between Go and TypeScript

## Tech Stack

### Backend

- **[Go](https://go.dev/doc/)** - High-performance backend server
- **[godotenv](https://github.com/joho/godotenv)** - Environment variable management

### Frontend

- **[React 19](https://react.dev/learn)** - UI framework
- **[TypeScript](https://www.typescriptlang.org/docs/)** - Type-safe frontend
- **[Vite](https://vite.dev/guide/)** - Fast build tool and dev server
- **[TanStack Router](https://tanstack.com/router/latest/docs/framework/react/overview)** - Type-safe routing
- **[TanStack Query](https://tanstack.com/query/latest/docs/framework/react/overview)** - Data fetching and state management
- **[Tailwind CSS](https://tailwindcss.com/docs/installation/using-vite)** - Utility-first CSS framework
- **[Zod](https://zod.dev/)** - Runtime type validation

### Tooling

- **[Moonrepo](https://moonrepo.dev/docs)** - Monorepo build system
- **[Bun](https://bun.com/docs)** - Package manager and runtime
- **[oxlint & oxfmt](https://oxc.rs/docs/guide/introduction.html)** - Fast TypeScript linter & Code formatter

## Project Structure

```
e2e-ts-go/
├── apps/
│   ├── server/          # Go backend server
│   │   ├── cmd/         # Application entry point
│   │   └── internal/    # Internal packages
│   │   └── moon.yml    # Server project configuration
│   ├── web/             # TypeScript frontend
│   │   └── src/
│   │       ├── components/
│   │       ├── lib/
│   │       └── routes/
│   │   └── moon.yml    # Web project configuration
├── .moon/
│   └── workspace.yml   # Moonrepo workspace configuration
└── package.json         # Root package.json for monorepo
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.25.5 or later
- [Bun](https://bun.sh/) 1.3.5 or later
- [Docker](https://www.docker.com/) (containerized development/deployment)

### Installation

```bash
# Install dependencies
bun install
```

### Environment Variables

Copy `.env.example` to `.env` in the root directory:

```bash
cp .env.example .env
```

Edit the `.env` file to set your environment variables as needed.

### Development

```bash
# Start Docker containers for required services (database, auth provider)
docker compose up -d --build

# Start both server and web app in development mode
bun dev
```

This will start:

- Go server (default port from `SERVER_PORT` env variable)
- Vite dev server for the frontend (default port from `VITE_WEB_PORT` env variable)

## Goals

- ✅ Go backend with HTTP server
- ✅ TypeScript frontend with modern React stack
- ✅ Monorepo setup with Moonrepo
- 🚧 End-to-end type safety between Go and TypeScript
- 🚧 Type generation from Go to TypeScript
- 🚧 RPC-like communication with type safety

## Inspiration

This project is inspired by [Better T Stack](https://www.better-t-stack.dev/), which provides a modern CLI for scaffolding end-to-end type-safe TypeScript projects. This implementation explores achieving similar type safety guarantees with a Go backend instead of a TypeScript/Node.js backend.

## License

See [LICENSE](./LICENSE) file for details.
