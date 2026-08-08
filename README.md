# Bank Credit Card Learning API

This repository is a learning workspace for a bank credit-card and customer-account domain. The current implementation blends a PostgreSQL relational schema, sample data under the SQL setup folder, API design documents, and a small Go HTTP service layer.

## Repository Purpose

The domain is modeled around customers, accounts, cards, transactions, payments, billing statements, fraud alerts, card limits, and rewards points. The source-of-truth database definitions are stored in [appdb_setup/appdb_table_ddl's.sql](appdb_setup/appdb_table_ddl's.sql), and the sample seed data lives in [appdb_setup/appdb_insert_statement.sql](appdb_setup/appdb_insert_statement.sql).

## Prerequisites

Before running the services or executing tests, make sure the following software is available locally:

- Go 1.22 or a compatible Go toolchain
- PostgreSQL server and client tooling such as psql
- Git for cloning and repository maintenance
- curl or another HTTP test tool for API smoke checks
- Optional: a local Swagger UI viewer or a static web server

## Setup Instructions

1. Clone the repository from your Git remote.
2. Change into the project root.
3. Download module dependencies:

```bash
go mod tidy
```

4. Make sure PostgreSQL is reachable locally and that the app database is created. The repository defaults expect the following connection values from [appdb_setup/appdb_connection_details.sql](appdb_setup/appdb_connection_details.sql):

- Host: localhost
- Port: 5432
- User: appuser
- Password: appuser123
- Database: appdb

Those defaults are read at runtime from environment variables:

- PGHOST
- PGPORT
- PGUSER
- PGPASSWORD
- PGDATABASE

If you want to override them, export them before starting each service. For example:

```bash
export PGHOST=localhost
export PGPORT=5432
export PGUSER=appuser
export PGPASSWORD=appuser123
export PGDATABASE=appdb
```

## Database Bootstrapping

The SQL scripts in this repo are the expected schema and sample-data input. Build the database from the DDL file and then load the seed rows:

```bash
psql -h localhost -p 5432 -U appuser -d appdb -f "appdb_setup/appdb_table_ddl's.sql"
psql -h localhost -p 5432 -U appuser -d appdb -f "appdb_setup/appdb_insert_statement.sql"
```

The DDL file is intentionally PostgreSQL-specific and should stay the starting point for future domain changes.

## Runtime Layout

The repository currently contains two lightweight Go microservices:

- Customer service: [services/customer-service/main.go](services/customer-service/main.go)
  - Default HTTP port: 8081
  - Main business area: customer identity and credit-card enrollment
  - Endpoints:
    - GET /health
    - GET /customers
    - GET /customers/{id}
    - POST /customers
    - PUT /customers/{id}

- Account service: [services/account-service/main.go](services/account-service/main.go)
  - Default HTTP port: 8082
  - Main business area: accounts, payments, and account-level transactions
  - Endpoints:
    - GET /health
    - GET /accounts
    - GET /accounts/{id}
    - GET /accounts/{id}/transactions
    - POST /payments

The shared database handle is created in [internal/db/postgres.go](internal/db/postgres.go) and reused by both service entrypoints.

## Running the Services Locally

From the repository root:

```bash
CUSTOMER_SERVICE_PORT=8081 go run ./services/customer-service
ACCOUNT_SERVICE_PORT=8082 go run ./services/account-service
```

Or, if you prefer the default values:

```bash
go run ./services/customer-service
go run ./services/account-service
```

You can probe the health endpoints with curl:

```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
```

## Customer API Example

The customer service supports the write/read API around the customers table. Example create request:

```bash
curl -X POST http://localhost:8081/customers \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Doe",
    "email": "jane.doe@example.com",
    "phone": "555-0111",
    "ssn": "987-65-4321",
    "customer_status": "active"
  }'
```

The response returns a customer payload with the generated `customer_id` and records the new row in the database.

## API Design and Contract Docs

The repository includes design artifacts:

- [docs/data-model.md](docs/data-model.md)
- [docs/design/HLD.md](docs/design/HLD.md)
- [docs/design/LLD.md](docs/design/LLD.md)
- [docs/design/API-contract.md](docs/design/API-contract.md)

The OpenAPI contract is in [api/openapi.yaml](api/openapi.yaml). There is also a static Swagger UI shell at [docs/swagger/index.html](docs/swagger/index.html) that loads that YAML.

To view Swagger locally through a web server:

```bash
python3 -m http.server 8000 --directory docs/swagger
```

Then open the browser URL:

http://localhost:8000/index.html

## Test Strategy

The workspace now carries a first-pass Godog/feature-style acceptance suite. The Gherkin scenario file is in [features/customer_api.feature](features/customer_api.feature), and the Go feature-binding test entrypoint is in [features/customer_api_test.go](features/customer_api_test.go).

Useful test commands:

```bash
go test ./features -run TestCustomerAPIFeature -v
go test ./services/customer-service
go test ./services/account-service
go test ./...
```

For a coverage-style verification run:

```bash
go test ./... -cover
```

## Notes for Extensibility

The codebase uses a schema-first approach. New services or API handlers should:

- reuse the existing domain vocabulary from the SQL model
- keep DB access centralized in [internal/db/postgres.go](internal/db/postgres.go)
- use table and column names that stay aligned with the DDL and sample data
- add feature-level BDD coverage whenever new REST behavior is introduced

## Current Status

The repository is a learning and bootstrap implementation for a small bank-card backend. Customer read/create/update flow is wired to PostgreSQL through the service layer, while account and payment analysis patterns are available as domain-first API entrypoints.

