# Tezos Delegation Service

A Go-based service that watches the Tezos blockchain for new blocks and delegation operations, storing delegation data in a PostgreSQL database and exposing it via a REST API.

## Features

- On start, it will get all delegations from tzkt and store them in the database
- Watches the Tezos blockchain for new blocks and delegations
- Stores delegation operations in a PostgreSQL database
- Exposes REST API endpoints to query delegations (optionally by year)
- Provides Prometheus-compatible metrics endpoint

## Quick Start

### 1. Prerequisites

- Docker & Docker Compose

### 2. Running with Docker Compose

```bash
docker compose up
```

This will start both the API and a PostgreSQL database.

### 3. Running Tests

```bash
make test
```

## Configuration

Configuration can be set via `config.yaml` (or `config-docker.yaml` for Docker):

```yaml
server:
  host: "localhost"
  port: 3000
  metricsPort: 3001

tzkt:
  url: "https://api.tzkt.io"

db:
  host: db
  port: 5432
  user: "postgres"
  password: "postgres"
  database: "delegations"
```

## API Endpoints

### Get All Delegations

```
GET /delegations
```

- Returns first 50 delegations in the database (ordered by year descending).

### Get Delegations by Year

```
GET /delegations/:year
```

- Returns delegations for the specified year.

#### Example:

```bash
curl http://localhost:3000/delegations
curl http://localhost:3000/delegations/2018
```

### Metrics

```
GET /metrics
```

- Prometheus metrics endpoint (default port 3001).

## Project Structure

- `main.go`: Entry point, wiring config, DB, watcher, and API
- `api/`: HTTP API and controllers
- `delegations_watcher/`: Watches Tezos chain and stores delegations
- `db/`: Database logic
- `httpclient/`: HTTP abstraction
- `middlewares/`: Gin middleware
- `types/`: Data types
- `utils/`: Utilities
