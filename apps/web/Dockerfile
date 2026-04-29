# Stage 1: Install dependencies
FROM oven/bun:1 AS deps

WORKDIR /app

COPY package.json apps/web/package.json bun.lock ./
RUN --mount=type=cache,target=/root/.bun/install/cache \
    bun install

# Stage 2: Build the application
FROM deps AS builder

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY apps/web/. .

ARG VITE_SERVER_URL=http://localhost:8080
ARG VITE_WEB_PORT=8081
ENV VITE_SERVER_URL=$VITE_SERVER_URL
ENV VITE_WEB_PORT=$VITE_WEB_PORT
ENV VITE_NODE_ENV=production

RUN bunx vite build

# Stage 3: Production runtime
FROM nginx:alpine AS release

RUN apk add --no-cache curl && rm -rf /var/cache/apk/*

COPY --from=builder /app/dist /usr/share/nginx/html
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 8081

CMD ["nginx", "-g", "daemon off;"]