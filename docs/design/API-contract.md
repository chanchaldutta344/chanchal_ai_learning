# API Contract

This repository exposes a small REST API layer over the existing bank credit-card domain tables.

## Customer service

Base URL: http://localhost:8081

### GET /health

Returns service health.

### GET /customers

Returns all customer rows.

### POST /customers

Creates a customer row using the customers table schema.

Request body:

```json
{
  "first_name": "Jane",
  "last_name": "Doe",
  "email": "jane.doe@example.com",
  "phone": "555-0111",
  "ssn": "987-65-4321",
  "customer_status": "active"
}
```

### GET /customers/{id}

Returns a single customer by customer_id.

### PUT /customers/{id}

Updates a single customer row.

## Account service

Base URL: http://localhost:8082

### GET /health

Returns service health.

### GET /accounts

Returns all account rows.

### GET /accounts/{id}

Returns one account row for account_id.

### GET /accounts/{id}/transactions

Returns all transaction rows for a given account_id.

### POST /payments

Creates a payment against the payments table.

Example body:

```json
{
  "account_id": 1,
  "payment_amount": 200.00,
  "payment_method": "Bank Transfer"
}
```

## Swagger / OpenAPI

The OpenAPI contract is stored in [api/openapi.yaml](api/openapi.yaml). It is designed to be viewable in a Swagger UI or future API documentation system.

Because there is no Swagger UI runtime pre-installed, a direct documentation link is not yet served by a local web server. The repository documents the YAML contract so it can be published into a Swagger UI or another docs site later.
