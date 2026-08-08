# Bank Credit Card Learning Agent Guide

This repository is a small learning workspace for a bank credit-card and customer-account domain. The current implementation assets are SQL DDL and sample insert scripts rather than a full web or API application.

## Project Shape

- Primary documentation entry point: [README.md](README.md)
- The database schema lives in [appdb_setup/appdb_table_ddl's.sql](appdb_setup/appdb_table_ddl's.sql)
- Sample bank/credit-card data is in [appdb_setup/appdb_insert_statement.sql](appdb_setup/appdb_insert_statement.sql)
- PostgreSQL connection environment is listed in [appdb_setup/appdb_connection_details.sql](appdb_setup/appdb_connection_details.sql)

## Business Domain

The data model is centered on a banking product offering:

- Customers hold profile and identity information.
- Credit cards belong to a customer and are issued by a bank or financial institution.
- Accounts link a customer to a card and carry credit-line, available credit, balance, and payment settings.
- Transactions record purchase activity against a card/account.
- Billing statements summarize balances and payment obligations.
- Payments represent payment activity applied to an account.
- Fraud alerts capture risk signals and case disposition.
- Rewards points capture earned/redeemed points associated with a transaction.
- Card limits express operational and risk controls.

The names in the schema are intentionally domain-oriented and should stay consistent when adding code or test fixtures.

## Database Conventions

- Use lower_snake_case names for tables and columns.
- Keep relationship keys explicit: `customer_id`, `card_id`, `account_id`, and `transaction_id` should stay consistent with the DDL.
- Prefer existing business vocabulary like `credit_cards`, `accounts`, `transactions`, `billing_statements`, `payments`, `fraud_alerts`, `rewards_points`, and `card_limits`.
- Preserve status-like fields such as `card_status`, `account_status`, `transaction_status`, `statement_status`, `payment_status`, `alert_status`, and `limit_status`.
- Treat monetary fields as numeric/decimal values rather than string or float placeholders.

## Data and Identity Rules

- Customer records should have unique email, phone, and SSN-style identifiers where appropriate.
- Credit card numbers are stored as TEXT and must remain unique.
- Account balance and credit capacity are derived from the account and card configuration, so future logic should not bypass those relationships.
- Fraud investigation artifacts are tied to a card and may optionally reference an account.
- Rewards point row-level logic should remain tied to a transaction when possible.

## Development Expectations

This workspace does not yet contain an application runtime, packaging metadata, or test harness. Agents should therefore:

- Prefer schema-first edits and SQL-safe reasoning.
- Reuse existing sample data patterns when adding example rows.
- Keep domain semantics intact when introducing new entities or reports.
- Avoid inventing UI or API layers unless other files are added to the repository.

## Working Guidance for AI Agents

When asked to modify or extend the codebase:

1. Check the schema before adding new tables or columns.
2. Reuse the existing domain vocabulary rather than creating parallel terminology.
3. Keep foreign-key relationships explicit.
4. When generating SQL, match the style of the DDL and insert files.
5. When writing tests or demo data, make the records readable and consistent with the sample values.

## Observed Wiring

The database environment points to PostgreSQL settings for a local app database:

- Host: `localhost`
- Port: `5432`
- User: `appuser`
- Password: `appuser123`
- Database: `appdb`

These settings should stay discoverable and should not be accidentally hard-coded into unrelated files.
