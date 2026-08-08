-- Insert Statements for Bank Credit Card Database
-- Sample data for testing and development

-- ============================================
-- INSERT INTO CUSTOMERS (4 records)
-- ============================================

INSERT INTO customers (customer_id, first_name, last_name, email, phone, ssn, date_of_birth, address, city, state, zip_code, country, customer_status)
VALUES 
(1, 'John', 'Smith', 'john.smith@email.com', '555-0101', '123-45-6789', '1985-03-15', '123 Main St', 'New York', 'NY', '10001', 'USA', 'active'),
(2, 'Sarah', 'Johnson', 'sarah.johnson@email.com', '555-0102', '234-56-7890', '1990-07-22', '456 Oak Ave', 'Los Angeles', 'CA', '90001', 'USA', 'active'),
(3, 'Michael', 'Davis', 'michael.davis@email.com', '555-0103', '345-67-8901', '1988-11-10', '789 Pine Rd', 'Chicago', 'IL', '60601', 'USA', 'active'),
(4, 'Emily', 'Wilson', 'emily.wilson@email.com', '555-0104', '456-78-9012', '1992-05-18', '321 Elm St', 'Houston', 'TX', '77001', 'USA', 'active');

-- ============================================
-- INSERT INTO CREDIT_CARDS (5 records)
-- ============================================

INSERT INTO credit_cards (card_id, customer_id, card_number, card_type, issuer_name, expiry_date, cvv, card_holder_name, card_status, issued_date)
VALUES
(1, 1, '4532015112830366', 'Visa', 'Chase Bank', '12/26', '123', 'John Smith', 'active', '2023-01-15'),
(2, 1, '4525233010103442', 'Mastercard', 'Bank of America', '08/25', '456', 'John Smith', 'active', '2022-06-10'),
(3, 2, '4516338506082832', 'Visa', 'Wells Fargo', '11/27', '789', 'Sarah Johnson', 'active', '2023-03-20'),
(4, 3, '4589041234567890', 'Mastercard', 'Citibank', '09/26', '321', 'Michael Davis', 'active', '2023-05-12'),
(5, 4, '4524007134432509', 'Visa', 'Capital One', '06/25', '654', 'Emily Wilson', 'active', '2022-09-05');

-- ============================================
-- INSERT INTO ACCOUNTS (5 records)
-- ============================================

INSERT INTO accounts (account_id, customer_id, card_id, credit_limit, available_credit, current_balance, minimum_payment, statement_close_day, payment_due_day, account_status, opened_date)
VALUES
(1, 1, 1, 15000.00, 12500.00, 2500.00, 75.00, 15, 10, 'active', '2023-01-15'),
(2, 1, 2, 10000.00, 8500.00, 1500.00, 45.00, 20, 15, 'active', '2022-06-10'),
(3, 2, 3, 20000.00, 18000.00, 2000.00, 60.00, 10, 5, 'active', '2023-03-20'),
(4, 3, 4, 12000.00, 10200.00, 1800.00, 54.00, 25, 20, 'active', '2023-05-12'),
(5, 4, 5, 18000.00, 16500.00, 1500.00, 45.00, 30, 25, 'active', '2022-09-05');

-- ============================================
-- INSERT INTO TRANSACTIONS (5 records)
-- ============================================

INSERT INTO transactions (transaction_id, account_id, card_id, merchant_name, merchant_category, transaction_amount, transaction_date, posting_date, transaction_type, description, reference_number, transaction_status)
VALUES
(1, 1, 1, 'Amazon.com', 'Electronics', 125.50, '2026-08-01 10:30:00', '2026-08-02 09:00:00', 'Purchase', 'Electronics Purchase', 'REF001', 'posted'),
(2, 1, 1, 'Starbucks', 'Restaurants', 5.75, '2026-08-02 08:15:00', '2026-08-02 09:30:00', 'Purchase', 'Coffee Purchase', 'REF002', 'posted'),
(3, 2, 2, 'Shell Gas Station', 'Gas', 55.00, '2026-08-03 14:45:00', '2026-08-04 09:00:00', 'Purchase', 'Fuel Purchase', 'REF003', 'posted'),
(4, 3, 3, 'Whole Foods Market', 'Groceries', 189.25, '2026-08-04 18:20:00', '2026-08-05 09:15:00', 'Purchase', 'Grocery Shopping', 'REF004', 'posted'),
(5, 4, 4, 'Delta Airlines', 'Travel', 450.00, '2026-08-05 09:00:00', '2026-08-05 10:00:00', 'Purchase', 'Flight Booking', 'REF005', 'posted');

-- ============================================
-- INSERT INTO PAYMENTS (3 records)
-- ============================================

INSERT INTO payments (payment_id, account_id, payment_amount, payment_date, payment_method, payment_status, reference_number, received_date)
VALUES
(1, 1, 1000.00, '2026-08-06 10:00:00', 'Bank Transfer', 'completed', 'PAY001', '2026-08-06 11:30:00'),
(2, 2, 750.00, '2026-08-07 14:30:00', 'Credit Transfer', 'completed', 'PAY002', '2026-08-07 15:45:00'),
(3, 3, 1200.00, '2026-08-08 09:15:00', 'Check', 'completed', 'PAY003', '2026-08-08 10:00:00');

-- ============================================
-- INSERT INTO BILLING_STATEMENTS (3 records)
-- ============================================

INSERT INTO billing_statements (statement_id, account_id, statement_date, due_date, previous_balance, current_charges, payments, current_balance, minimum_payment, interest_charged, interest_rate, statement_status)
VALUES
(1, 1, '2026-08-01', '2026-08-25', 1500.00, 500.00, 250.00, 1750.00, 52.50, 15.00, 18.99, 'generated'),
(2, 2, '2026-08-01', '2026-08-26', 800.00, 350.00, 100.00, 1050.00, 31.50, 8.50, 18.99, 'generated'),
(3, 3, '2026-08-01', '2026-08-25', 1200.00, 600.00, 400.00, 1400.00, 42.00, 12.00, 18.99, 'generated');

-- ============================================
-- INSERT INTO FRAUD_ALERTS (2 records)
-- ============================================

INSERT INTO fraud_alerts (alert_id, card_id, account_id, alert_type, alert_description, alert_date, alert_status, resolution_notes, resolved_date)
VALUES
(1, 1, 1, 'Unusual Activity', 'Multiple transactions in different locations within 1 hour', '2026-08-05 15:30:00', 'resolved', 'Verified legitimate transactions - customer traveling', '2026-08-05 16:00:00'),
(2, 3, 3, 'Large Transaction', 'Transaction amount exceeds typical spending pattern', '2026-08-04 18:45:00', 'resolved', 'Confirmed airline purchase by customer', '2026-08-04 19:30:00');

-- ============================================
-- INSERT INTO REWARDS_POINTS (4 records)
-- ============================================

INSERT INTO rewards_points (rewards_id, account_id, card_id, points_earned, points_redeemed, total_points, points_expiry_date, transaction_id, earned_date)
VALUES
(1, 1, 1, 1.26, 0.00, 1.26, '2027-08-01', 1, '2026-08-01 10:30:00'),
(2, 1, 1, 0.06, 0.00, 1.32, '2027-08-02', 2, '2026-08-02 08:15:00'),
(3, 2, 2, 0.55, 0.00, 0.55, '2027-08-03', 3, '2026-08-03 14:45:00'),
(4, 3, 3, 1.89, 0.00, 1.89, '2027-08-04', 4, '2026-08-04 18:20:00');

-- ============================================
-- INSERT INTO CARD_LIMITS (3 records)
-- ============================================

INSERT INTO card_limits (limit_id, card_id, account_id, daily_transaction_limit, daily_withdrawal_limit, monthly_transaction_count_limit, international_limit, limit_effective_date, limit_status)
VALUES
(1, 1, 1, 5000.00, 2000.00, 100, 3000.00, '2023-01-15', 'active'),
(2, 3, 3, 7500.00, 3000.00, 150, 5000.00, '2023-03-20', 'active'),
(3, 4, 4, 6000.00, 2500.00, 120, 4000.00, '2023-05-12', 'active');
