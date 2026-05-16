# Architecture

- **Backend**: Go server with HTTP handlers
- **Frontend**: TypeScript/React with Vite, TanStack Router, and TanStack Query
- **Monorepo**: Managed with Moonrepo for efficient builds and development
- **Type Safety**: End-to-end type safety between Go and TypeScript

## Tech Stack

### Backend

- **[Go](https://go.dev/doc/)** - High-performance backend server
- **[Authula](https://authula.vercel.app/docs/)** - Authentication
- **[Bun](https://bun.uptrace.dev/)** - ORM
- **[Chi](https://go-chi.io/)** - Router
- **[Faker](https://github.com/jaswdr/faker)** - Fake data generation
- **[Swagger](https://github.com/swaggo/swag) & [Scalar](https://scalar.com/)** - API documentation

### Frontend

- **[React 19](https://react.dev/learn)** - UI framework
- **[TypeScript](https://www.typescriptlang.org/docs/)** - Type-safe frontend
- **[Vite](https://vite.dev/guide/)** - Fast build tool and dev server
- **[TanStack Router](https://tanstack.com/router/latest/docs/framework/react/overview)** - Type-safe routing
- **[TanStack Query](https://tanstack.com/query/latest/docs/framework/react/overview)** - Data fetching and state management
- **[Tailwind CSS](https://tailwindcss.com/docs/installation/using-vite)** - Utility-first CSS framework
- **[Zod](https://zod.dev/)** - Runtime type validation
- **[oRPC](https://orpc.unnoq.com/)** - End-to-end type-safe API
- **[HeyApi](https://heyapi.dev/)** - OpenAPI client generation

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

- [Go](https://go.dev/) 1.26.2 or later
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

# Seed the database
bun run db-seed
```

This will start:

- Go server (default port from `SERVER_PORT` env variable)
- Vite dev server for the frontend (default port from `VITE_WEB_PORT` env variable)

## License

See [LICENSE](./LICENSE) file for details.
