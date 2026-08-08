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
    "bankcard-learning/internal/httpx"
)

// Customer models the customer-facing identity entity used by the customer
// service API. The struct mirrors the customers table shape in the DDL.
type Customer struct {
    CustomerID int    `json:"customer_id"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Email      string `json:"email"`
    Phone      string `json:"phone"`
    SSN        string `json:"ssn"`
    Status     string `json:"customer_status"`
}

// CreditCard is kept as a lightweight API model and is not yet exposed as a
// working REST resource in this first implementation.
type CreditCard struct {
    CardID       int    `json:"card_id"`
    CustomerID   int    `json:"customer_id"`
    CardNumber   string `json:"card_number"`
    CardType     string `json:"card_type"`
    IssuerName   string `json:"issuer_name"`
    ExpiryDate   string `json:"expiry_date"`
    CardStatus   string `json:"card_status"`
    CardHolder   string `json:"card_holder_name"`
    IssuedDate   string `json:"issued_date"`
}

// main boots the customer service and binds its HTTP routes to the database
// handle created from the shared Postgres configuration.
func main() {
    db, err := db.InitDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // A service-level router keeps the resource ownership isolated and keeps
    // route registration easy to extend as the domain grows. The router is a
    // tiny path switcher rather than a ServeMux shorthand because the latter can
    // interpret a /customers/ request as a collection-level match when the
    // patterns are registered in tension with one another.
    mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch {
        case r.URL.Path == "/health":
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            _, _ = w.Write([]byte(`{"status":"ok"}`))
        case r.URL.Path == "/customers":
            switch r.Method {
            case http.MethodGet:
                listCustomers(w, db)
            case http.MethodPost:
                createCustomer(w, r, db)
            default:
                w.WriteHeader(http.StatusMethodNotAllowed)
            }
        case strings.HasPrefix(r.URL.Path, "/customers/"):
            if r.URL.Path == "/customers/" {
                w.WriteHeader(http.StatusNotFound)
                return
            }

            pathID := strings.TrimPrefix(r.URL.Path, "/customers/")
            customerID, err := strconv.Atoi(pathID)
            if err != nil {
                http.Error(w, "invalid customer id", http.StatusBadRequest)
                return
            }

            switch r.Method {
            case http.MethodGet:
                getCustomer(w, db, customerID)
            case http.MethodPut:
                updateCustomer(w, r, db, customerID)
            default:
                w.WriteHeader(http.StatusMethodNotAllowed)
            }
        default:
            http.NotFound(w, r)
        }
    })

    port := os.Getenv("CUSTOMER_SERVICE_PORT")
    if port == "" {
        port = "8081"
    }

    fmt.Println("Customer service listening on :" + port)
    log.Fatal(http.ListenAndServe(":"+port, mux))
}

// listCustomers reads the customers table and returns a JSON array. It is the
// read-model support for the customer collection resource.
func listCustomers(w http.ResponseWriter, db *sql.DB) {
    rows, err := db.Query("SELECT customer_id, first_name, last_name, email, phone, ssn, customer_status FROM customers")
    if err != nil {
        httpx.WriteServerError(w, "list customers", err)
        return
    }
    defer rows.Close()

    customers := []Customer{}
    for rows.Next() {
        var c Customer
        if err := rows.Scan(&c.CustomerID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.SSN, &c.Status); err != nil {
            httpx.WriteServerError(w, "scan customer row", err)
            return
        }
        customers = append(customers, c)
    }
    if err := rows.Err(); err != nil {
        httpx.WriteServerError(w, "iterate customers", err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(customers)
}

// getCustomer performs the id-level lookup of a customer row.
func getCustomer(w http.ResponseWriter, db *sql.DB, customerID int) {
    var c Customer
    err := db.QueryRow("SELECT customer_id, first_name, last_name, email, phone, ssn, customer_status FROM customers WHERE customer_id = $1", customerID).Scan(
        &c.CustomerID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.SSN, &c.Status,
    )
    if err == sql.ErrNoRows {
        http.Error(w, "customer not found", http.StatusNotFound)
        return
    }
    if err != nil {
        httpx.WriteServerError(w, "get customer", err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(c)
}

// createCustomer handles the POST /customers write path. It derives the next
// customer key from the max PK and writes a full row into the customers table.
func createCustomer(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    var c Customer
    if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // customers.customer_id has no sequence in the DDL, so the key is derived
    // inside the INSERT itself. A separate SELECT MAX followed by an INSERT lets
    // two concurrent creates pick the same key and fail on customers_pkey.
    row := db.QueryRow(
        `INSERT INTO customers (customer_id, first_name, last_name, email, phone, ssn, customer_status)
         SELECT COALESCE(MAX(customer_id), 0) + 1, $1, $2, $3, $4, $5, $6 FROM customers
         RETURNING customer_id`,
        c.FirstName,
        c.LastName,
        c.Email,
        c.Phone,
        c.SSN,
        c.Status,
    )

    if err := row.Scan(&c.CustomerID); err != nil {
        if httpx.IsUniqueViolation(err) {
            http.Error(w, "a customer with the same unique identity already exists", http.StatusConflict)
            return
        }
        httpx.WriteServerError(w, "create customer", err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(c)
}

// updateCustomer is the PUT handler for /customers/{id}. It updates the row
// that is already confirmed by the customer_id in the path.
func updateCustomer(w http.ResponseWriter, r *http.Request, db *sql.DB, customerID int) {
    var update Customer
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if update.CustomerID != 0 && update.CustomerID != customerID {
        http.Error(w, "customer_id in payload does not match path", http.StatusBadRequest)
        return
    }

    update.CustomerID = customerID

    result, err := db.Exec(
        "UPDATE customers SET first_name = $1, last_name = $2, email = $3, phone = $4, ssn = $5, customer_status = $6 WHERE customer_id = $7",
        update.FirstName,
        update.LastName,
        update.Email,
        update.Phone,
        update.SSN,
        update.Status,
        update.CustomerID,
    )
    if err != nil {
        if httpx.IsUniqueViolation(err) {
            http.Error(w, "a customer with the same unique identity already exists", http.StatusConflict)
            return
        }
        httpx.WriteServerError(w, "update customer", err)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        httpx.WriteServerError(w, "read update result", err)
        return
    }
    if rowsAffected == 0 {
        http.Error(w, "customer not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(update)
}
