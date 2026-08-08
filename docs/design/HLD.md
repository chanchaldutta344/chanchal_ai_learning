# High-Level Design (HLD) — Bank Credit Card REST API

## 1. Purpose

This repository is being developed as a small bank credit-card learning platform with a relational data model and simple REST APIs for transactional business operations. The initial scope is customer creation, read access, customer updates, account lookup, account transaction history, and payment recording.

The API design is intentionally small and service-oriented:

- Customer service handles customer identity and basic card-holder lifecycle operations.
- Account service handles account-level financial exposure, transaction history, and payment operations.

## 2. Business Domain

The supported domain is a banking and credit-card domain covering the following business objects:

- Customer
- Credit card
- Account
- Transaction
- Billing statement
- Payment
- Fraud alert
- Rewards points
- Card limits

The HLD is centered around a customer-to-account-to-transaction financial flow.

## 3. Current Architecture

### 3.1 Service Decomposition

| Service | Responsibility | Port |
|---|---|---|
| Customer service | Customer registration, lookup, update | 8081 |
| Account service | Account lookup, transaction history, payment posting | 8082 |

### 3.2 Data Store

The persistent store is PostgreSQL and the source-of-truth schema is defined in the SQL DDL under [appdb_setup/appdb_table_ddl's.sql](appdb_setup/appdb_table_ddl's.sql). The runtime uses the env values from [appdb_setup/appdb_connection_details.sql](appdb_setup/appdb_connection_details.sql):

- Host: localhost
- Port: 5432
- User: appuser
- Password: appuser123
- Database: appdb

### 3.3 API Style

The system uses REST over HTTP/JSON. A simple Go HTTP server is used as the first implementation. Each service is responsible for a subset of domain resources and writes directly to PostgreSQL using the sql package.

## 4. Functional Scope

### Customer Service

The customer service implements REST operations around the `customers` table:

- Create a customer
- List customers
- Get a customer by ID
- Update a customer by ID

### Account Service

The account service implements operations around the `accounts` and `payments` table domain:

- List accounts
- Get an account by ID
- Get account transaction history
- Create a payment record

## 5. API Endpoints

### Customer Service Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /health | Health check |
| GET | /customers | List customer rows |
| POST | /customers | Create customer |
| GET | /customers/{id} | Get single customer |
| PUT | /customers/{id} | Update a single customer |

### Account Service Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | /health | Health check |
| GET | /accounts | List accounts |
| GET | /accounts/{id} | Get account details |
| GET | /accounts/{id}/transactions | List transactions for an account |
| POST | /payments | Create payment record |

## 6. Business Rules

- Customer creation must supply the required `customer_id` PK value through SQL generation logic.
- Customer email, phone, and SSN should be unique in future API validation.
- Account reads should be constrained to a valid existing customer/card relationship.
- Payment records should reference a valid account in the `accounts` table.
- Transaction reads should be account-scoped.

## 7. Non-Functional Requirements

### Reliability

- The API should return clear 4xx or 5xx responses when business or schema constraints fail.
- The service should not invent invalid IDs or break the existing PK/FK flow.

### Security

- The initial implementation uses local dev credentials and is not production hardened.
- Future work should add authentication, role-based access, secrets management, and request validation.

### Observability

- Basic health endpoints should be available for service health checks.
- Logs should be surfaced for database exceptions and HTTP failures.

## 8. Future Roadmap

The immediate next extensions could include:

1. Credit-card creation service or card enrollment endpoint.
2. Billing-statement generation and retrieval endpoint.
3. Fraud alert review endpoints.
4. Rewards points lookup endpoint.
5. API gateway / BFF layer.
6. Database migrations using SQL schema versioning.
7. Test coverage for happy-path and negative SQL error paths.
