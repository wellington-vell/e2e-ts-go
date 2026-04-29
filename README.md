# Architecture

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
- **[Oxc and Golangci-lint](https://oxc.rs/docs/guide/introduction.html)** - Linters & Code formatters

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
- [Golangci Lint](https://golangci-lint.run/docs/welcome/install/local/) linter for golang language

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
# Start Docker dev containers
bun run docker-dev

# Start both server and web app in development mode
bun dev
```

This will start:

- Go server (default port from `SERVER_PORT` env variable)
- Vite dev server for the frontend (default port from `VITE_WEB_PORT` env variable)

### Database Migrations

Database migrations are managed with [goose](https://github.com/pressly/goose). Make sure `DATABASE_URL` is set in your `.env` file.

```bash
# Create a new migration
bun run migrate-create -- add_users sql

# Run all pending migrations
bun run migrate-up

# Roll back the last migration
bun run migrate-down

# Migrate to a specific version
bun run migrate-up-to -- 20260425000000

# Roll back to version 0 (all migrations)
bun run migrate-down-to -- 0
```

## License

See [LICENSE](./LICENSE) file for details.
