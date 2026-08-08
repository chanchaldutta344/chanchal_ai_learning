CREATE TABLE customers (
  customer_id INTEGER PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  phone TEXT,
  ssn TEXT UNIQUE NOT NULL,
  date_of_birth DATE,
  address TEXT,
  city TEXT,
  state TEXT,
  zip_code TEXT,
  country TEXT,
  customer_status TEXT DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE credit_cards (
  card_id INTEGER PRIMARY KEY,
  customer_id INTEGER NOT NULL,
  card_number TEXT UNIQUE NOT NULL,
  card_type TEXT NOT NULL,
  issuer_name TEXT NOT NULL,
  expiry_date TEXT NOT NULL,
  cvv TEXT NOT NULL,
  card_holder_name TEXT NOT NULL,
  card_status TEXT DEFAULT 'active',
  issued_date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);

CREATE TABLE accounts (
  account_id INTEGER PRIMARY KEY,
  customer_id INTEGER NOT NULL,
  card_id INTEGER NOT NULL,
  credit_limit DECIMAL(12, 2) NOT NULL,
  available_credit DECIMAL(12, 2) NOT NULL,
  current_balance DECIMAL(12, 2) DEFAULT 0,
  minimum_payment DECIMAL(12, 2),
  statement_close_day INTEGER,
  payment_due_day INTEGER,
  account_status TEXT DEFAULT 'active',
  opened_date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
  FOREIGN KEY (card_id) REFERENCES credit_cards(card_id)
);

CREATE TABLE transactions (
  transaction_id INTEGER PRIMARY KEY,
  account_id INTEGER NOT NULL,
  card_id INTEGER NOT NULL,
  merchant_name TEXT NOT NULL,
  merchant_category TEXT,
  transaction_amount DECIMAL(12, 2) NOT NULL,
  transaction_date TIMESTAMP NOT NULL,
  posting_date TIMESTAMP,
  transaction_type TEXT NOT NULL,
  description TEXT,
  reference_number TEXT UNIQUE,
  transaction_status TEXT DEFAULT 'posted',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id),
  FOREIGN KEY (card_id) REFERENCES credit_cards(card_id)
);

CREATE TABLE billing_statements (
  statement_id INTEGER PRIMARY KEY,
  account_id INTEGER NOT NULL,
  statement_date DATE NOT NULL,
  due_date DATE NOT NULL,
  previous_balance DECIMAL(12, 2) NOT NULL,
  current_charges DECIMAL(12, 2) NOT NULL,
  payments DECIMAL(12, 2) DEFAULT 0,
  current_balance DECIMAL(12, 2) NOT NULL,
  minimum_payment DECIMAL(12, 2) NOT NULL,
  interest_charged DECIMAL(12, 2) DEFAULT 0,
  interest_rate DECIMAL(5, 2),
  statement_status TEXT DEFAULT 'generated',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE payments (
  payment_id INTEGER PRIMARY KEY,
  account_id INTEGER NOT NULL,
  payment_amount DECIMAL(12, 2) NOT NULL,
  payment_date TIMESTAMP NOT NULL,
  payment_method TEXT NOT NULL,
  payment_status TEXT DEFAULT 'completed',
  reference_number TEXT UNIQUE,
  received_date TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE fraud_alerts (
  alert_id INTEGER PRIMARY KEY,
  card_id INTEGER NOT NULL,
  account_id INTEGER,
  alert_type TEXT NOT NULL,
  alert_description TEXT,
  alert_date TIMESTAMP NOT NULL,
  alert_status TEXT DEFAULT 'open',
  resolution_notes TEXT,
  resolved_date TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (card_id) REFERENCES credit_cards(card_id),
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE TABLE rewards_points (
  rewards_id INTEGER PRIMARY KEY,
  account_id INTEGER NOT NULL,
  card_id INTEGER NOT NULL,
  points_earned DECIMAL(10, 2) NOT NULL,
  points_redeemed DECIMAL(10, 2) DEFAULT 0,
  total_points DECIMAL(10, 2),
  points_expiry_date DATE,
  transaction_id INTEGER,
  earned_date TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (account_id) REFERENCES accounts(account_id),
  FOREIGN KEY (card_id) REFERENCES credit_cards(card_id),
  FOREIGN KEY (transaction_id) REFERENCES transactions(transaction_id)
);

CREATE TABLE card_limits (
  limit_id INTEGER PRIMARY KEY,
  card_id INTEGER NOT NULL,
  account_id INTEGER NOT NULL,
  daily_transaction_limit DECIMAL(12, 2),
  daily_withdrawal_limit DECIMAL(12, 2),
  monthly_transaction_count_limit INTEGER,
  international_limit DECIMAL(12, 2),
  limit_effective_date DATE,
  limit_status TEXT DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (card_id) REFERENCES credit_cards(card_id),
  FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);