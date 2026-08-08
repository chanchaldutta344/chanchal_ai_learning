package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "bankcard-learning/internal/db"
)

// Account models the account table fields for the account-oriented service.
// The service uses it as the JSON contract for GET /accounts and GET /accounts/:id.
type Account struct {
    AccountID          int     `json:"account_id"`
    CustomerID         int     `json:"customer_id"`
    CardID             int     `json:"card_id"`
    CreditLimit        float64 `json:"credit_limit"`
    AvailableCredit    float64 `json:"available_credit"`
    CurrentBalance     float64 `json:"current_balance"`
    MinimumPayment     float64 `json:"minimum_payment"`
    StatementCloseDay  int     `json:"statement_close_day"`
    PaymentDueDay      int     `json:"payment_due_day"`
    AccountStatus      string  `json:"account_status"`
    OpenedDate         string  `json:"opened_date"`
}

// Transaction models the transaction fact table and drives the account-level
// statement and account-history read flow.
type Transaction struct {
    TransactionID    int     `json:"transaction_id"`
    AccountID        int     `json:"account_id"`
    CardID           int     `json:"card_id"`
    MerchantName     string  `json:"merchant_name"`
    MerchantCategory string  `json:"merchant_category"`
    TransactionAmount float64 `json:"transaction_amount"`
    TransactionDate  string  `json:"transaction_date"`
    PostingDate      string  `json:"posting_date"`
    TransactionType  string  `json:"transaction_type"`
    Description      string  `json:"description"`
    ReferenceNumber  string  `json:"reference_number"`
    TransactionStatus string `json:"transaction_status"`
}

// PaymentRequest captures the minimal payload required to create a payment row.
type PaymentRequest struct {
    AccountID int     `json:"account_id"`
    PaymentAmount float64 `json:"payment_amount"`
    PaymentMethod string `json:"payment_method"`
}

// main instantiates the account service and exposes account, transaction, and
// payment APIs over the same PostgreSQL-backed schema.
func main() {
    db, err := db.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // The router intentionally separates read-only list/detail routes from the
    // POST payment side-effecting resource. It can grow without changing DB shape.
    mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch {
        case r.URL.Path == "/health":
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            _, _ = w.Write([]byte(`{"status":"ok"}`))
        case r.URL.Path == "/accounts":
            switch r.Method {
            case http.MethodGet:
                listAccounts(w, db)
            default:
                w.WriteHeader(http.StatusMethodNotAllowed)
            }
        case strings.HasPrefix(r.URL.Path, "/accounts/"):
            if r.Method != http.MethodGet {
                w.WriteHeader(http.StatusMethodNotAllowed)
                return
            }

            path := r.URL.Path
            if path == "/accounts/" {
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            suffix := strings.TrimPrefix(path, "/accounts/")
            segments := strings.Split(suffix, "/")
            if len(segments) == 1 {
                accountID, err := strconv.Atoi(segments[0])
                if err != nil {
                    http.Error(w, "invalid account id", http.StatusBadRequest)
                    return
                }
                getAccount(w, db, accountID)
                return
            }

            if len(segments) == 2 && segments[1] == "transactions" {
                accountID, err := strconv.Atoi(segments[0])
                if err != nil {
                    http.Error(w, "invalid account id", http.StatusBadRequest)
                    return
                }
                listTransactionsForAccount(w, db, accountID)
                return
            }

            w.WriteHeader(http.StatusNotFound)
        case r.URL.Path == "/payments":
            switch r.Method {
            case http.MethodPost:
                createPayment(w, r, db)
            default:
                w.WriteHeader(http.StatusMethodNotAllowed)
            }
        default:
            http.NotFound(w, r)
        }
    })

    port := os.Getenv("ACCOUNT_SERVICE_PORT")
    if port == "" {
        port = "8082"
    }

    fmt.Println("Account service listening on :" + port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}

// listAccounts queries the account table and returns the account collection as JSON.
func listAccounts(w http.ResponseWriter, db *sql.DB) {
    rows, err := db.Query("SELECT account_id, customer_id, card_id, credit_limit, available_credit, current_balance, minimum_payment, statement_close_day, payment_due_day, account_status, opened_date FROM accounts")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    accounts := []Account{}
    for rows.Next() {
        var a Account
        if err := rows.Scan(&a.AccountID, &a.CustomerID, &a.CardID, &a.CreditLimit, &a.AvailableCredit, &a.CurrentBalance, &a.MinimumPayment, &a.StatementCloseDay, &a.PaymentDueDay, &a.AccountStatus, &a.OpenedDate); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        accounts = append(accounts, a)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(accounts)
}

// getAccount reads a single account by account_id.
func getAccount(w http.ResponseWriter, db *sql.DB, accountID int) {
    var a Account
    err := db.QueryRow("SELECT account_id, customer_id, card_id, credit_limit, available_credit, current_balance, minimum_payment, statement_close_day, payment_due_day, account_status, opened_date FROM accounts WHERE account_id = $1", accountID).Scan(
        &a.AccountID, &a.CustomerID, &a.CardID, &a.CreditLimit, &a.AvailableCredit, &a.CurrentBalance, &a.MinimumPayment, &a.StatementCloseDay, &a.PaymentDueDay, &a.AccountStatus, &a.OpenedDate,
    )
    if err == sql.ErrNoRows {
        http.Error(w, "account not found", http.StatusNotFound)
        return
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(a)
}

// listTransactionsForAccount is the transaction-reporting side of the account
// service and returns all related transaction rows for a known account.
func listTransactionsForAccount(w http.ResponseWriter, db *sql.DB, accountID int) {
    rows, err := db.Query("SELECT transaction_id, account_id, card_id, merchant_name, merchant_category, transaction_amount, transaction_date, posting_date, transaction_type, description, reference_number, transaction_status FROM transactions WHERE account_id = $1", accountID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    txns := []Transaction{}
    for rows.Next() {
        var t Transaction
        if err := rows.Scan(&t.TransactionID, &t.AccountID, &t.CardID, &t.MerchantName, &t.MerchantCategory, &t.TransactionAmount, &t.TransactionDate, &t.PostingDate, &t.TransactionType, &t.Description, &t.ReferenceNumber, &t.TransactionStatus); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        txns = append(txns, t)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(txns)
}

// createPayment inserts a row into the payments table. This is the current
// side-effecting API for account-level settlement operations.
func createPayment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var req PaymentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var paymentID int
    err := db.QueryRow(
        "INSERT INTO payments (account_id, payment_amount, payment_date, payment_method, payment_status, reference_number, received_date) VALUES ($1, $2, CURRENT_TIMESTAMP, $3, 'completed', 'PAY-' || nextval('payments_payment_id_seq'), CURRENT_TIMESTAMP) RETURNING payment_id",
        req.AccountID,
        req.PaymentAmount,
        req.PaymentMethod,
    ).Scan(&paymentID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]any{"payment_id": paymentID, "status": "created"})
}
