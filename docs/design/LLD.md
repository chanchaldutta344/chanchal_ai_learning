# Low-Level Design (LLD) — Bank Credit Card REST API

## 1. Scope

This low-level design describes the Go HTTP services currently present in the repository and the database operations each service performs.

## 2. Service Implementation Details

### 2.1 Customer Service

Implementation file:

- [services/customer-service/main.go](services/customer-service/main.go)

Package:

- `main`

HTTP router:

- `http.NewServeMux()`

Dependencies:

- `database/sql`
- PostgreSQL driver: `github.com/lib/pq`
- `encoding/json`
- standard `net/http`

#### Endpoints

GET /health

- Response: JSON status object
- Example response:

```json
{"status":"ok"}
```

GET /customers

- SQL: `SELECT customer_id, first_name, last_name, email, phone, ssn, customer_status FROM customers`
- Response shape: list of customer JSON structures.

POST /customers

- Request body expects customer fields.
- SQL insert performs:

```sql
INSERT INTO customers (
  customer_id,
  first_name,
  last_name,
  email,
  phone,
  ssn,
  customer_status
) VALUES ($1, $2, $3, $4, $5, $6, $7)
```

- Primary-key assignment logic: `SELECT COALESCE(MAX(customer_id), 0) + 1 FROM customers`
- Response: 201 Created with created customer payload.

GET /customers/{id}

- SQL: `SELECT ... FROM customers WHERE customer_id = $1`
- Response: one customer JSON structure.

PUT /customers/{id}

- Request body expects a customer payload.
- SQL update performs:

```sql
UPDATE customers
SET first_name = $1,
    last_name = $2,
    email = $3,
    phone = $4,
    ssn = $5,
    customer_status = $6
WHERE customer_id = $7
```

- Response: 200 OK with updated customer payload.

### 2.2 Account Service

Implementation file:

- [services/account-service/main.go](services/account-service/main.go)

Package:

- `main`

HTTP router:

- `http.NewServeMux()`

Dependencies:

- `database/sql`
- PostgreSQL driver: `github.com/lib/pq`
- `encoding/json`
- `net/http`

#### Endpoints

GET /health

- Response: JSON status object.

GET /accounts

- SQL: list rows from `accounts` table.
- Response: list of account payloads.

GET /accounts/{id}

- SQL: account detail by `account_id`.
- Response: account JSON.

GET /accounts/{id}/transactions

- SQL:

```sql
SELECT transaction_id, account_id, card_id, merchant_name, merchant_category, transaction_amount, transaction_date, posting_date, transaction_type, description, reference_number, transaction_status
FROM transactions
WHERE account_id = $1
```

- Response: array of transaction records.

POST /payments

- Request body expects an account id, a positive payment amount with at most two
  decimal places, and a supported payment method. The referenced account must
  already exist.
- `payments.payment_id` has no sequence in the DDL, so the key and its derived
  reference number are generated inside the insert statement itself:

```sql
INSERT INTO payments (
  payment_id,
  account_id,
  payment_amount,
  payment_date,
  payment_method,
  payment_status,
  reference_number,
  received_date
)
SELECT next_id, $1, $2, CURRENT_TIMESTAMP, $3, 'completed', 'PAY-' || next_id, CURRENT_TIMESTAMP
FROM (SELECT COALESCE(MAX(payment_id), 0) + 1 AS next_id FROM payments) AS candidate
RETURNING payment_id
```

- Response: 201 Created with a status payload.

## 3. Data Structures

### Customer

```go
type Customer struct {
    CustomerID int    `json:"customer_id"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    SSN        string `json:"ssn"`
    Status     string `json:"customer_status"`
}
```

### Account

```go
type Account struct {
    AccountID         int          `json:"account_id"`
    CustomerID        int          `json:"customer_id"`
    CardID            int          `json:"card_id"`
    CreditLimit       money.Amount `json:"credit_limit"`
    AvailableCredit   money.Amount `json:"available_credit"`
    CurrentBalance    money.Amount `json:"current_balance"`
    MinimumPayment    money.Amount `json:"minimum_payment"`
    StatementCloseDay int          `json:"statement_close_day"`
    PaymentDueDay     int          `json:"payment_due_day"`
    AccountStatus     string       `json:"account_status"`
    OpenedDate        string       `json:"opened_date"`
}
```

Monetary fields use `money.Amount` (an alias for `shopspring/decimal.Decimal`)
so the `DECIMAL(12, 2)` columns round-trip exactly instead of through binary
floating point. They are still encoded as JSON numbers.

### Transaction

```go
type Transaction struct {
    TransactionID     int          `json:"transaction_id"`
    AccountID         int          `json:"account_id"`
    CardID            int          `json:"card_id"`
    MerchantName      string       `json:"merchant_name"`
    MerchantCategory  string       `json:"merchant_category"`
    TransactionAmount money.Amount `json:"transaction_amount"`
    TransactionDate   string       `json:"transaction_date"`
    PostingDate       string       `json:"posting_date"`
    TransactionType   string       `json:"transaction_type"`
    Description       string       `json:"description"`
    ReferenceNumber   string       `json:"reference_number"`
    TransactionStatus string       `json:"transaction_status"`
}
```

## 4. Database Connection Layer

Both services call a local `initDB()` function that reads environment variables from the runtime process:

- `PGHOST`
- `PGPORT`
- `PGUSER`
- `PGPASSWORD`
- `PGDATABASE`

Host, port, user, and database name fall back to the values in the SQL
connection settings file. `PGPASSWORD` has no default: `InitDB` returns an error
when it is unset so a deployment can never fall back to a well-known password.

```text
PGHOST=localhost
PGPORT=5432
PGUSER=appuser
PGDATABASE=appdb
```

The DSN string used is:

```go
host=%s port=%s user=%s password=%s dbname=%s sslmode=disable
```

## 5. Failure Handling

The services map DB and HTTP errors to HTTP status codes:

- 400 Bad Request for malformed input and bad ID parameters.
- 404 Not Found for missing customer/account rows.
- 405 Method Not Allowed for unsupported HTTP verbs.
- 500 Internal Server Error for DB failures. The driver message is logged
  server-side only; the response body is a generic `internal server error` so
  schema, constraint, and row details never reach the client.

## 6. Testability

The current implementation is service-level and has no dedicated Go test suite. The repository uses `go test ./...` as the compile-level verification command. Future work should add table-driven API tests and database integration tests.

## 7. Deployment Notes

The services are intended to run locally on separate ports. In a future environment they can be deployed as containers behind an API gateway or service mesh. The current design is a good starting point for a synchronous REST-based integration style.
