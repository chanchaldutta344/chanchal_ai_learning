package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"

    "bankcard-learning/internal/db"
)

func TestCustomerCRUDLifecycle(t *testing.T) {
    dbConn, err := db.InitDB()
    if err != nil {
        t.Fatalf("failed to connect to postgres: %v", err)
    }
    defer dbConn.Close()

    mux := http.NewServeMux()
    mux.HandleFunc("/customers", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            listCustomers(w, dbConn)
        case http.MethodPost:
            createCustomer(w, r, dbConn)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })
    mux.HandleFunc("/customers/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/customers/" {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        pathID := r.URL.Path[len("/customers/"):]
        customerID, err := strconv.Atoi(pathID)
        if err != nil {
            http.Error(w, "invalid customer id", http.StatusBadRequest)
            return
        }
        switch r.Method {
        case http.MethodGet:
            getCustomer(w, dbConn, customerID)
        case http.MethodPut:
            updateCustomer(w, r, dbConn, customerID)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    server := httptest.NewServer(mux)
    defer server.Close()

    payload := `{
        "first_name":"Coverage",
        "last_name":"User",
        "email":"coverage.user.a@example.com",
        "phone":"555-0999",
        "ssn":"111-22-3333",
        "customer_status":"active"
    }`

    resp, err := http.Post(server.URL+"/customers", "application/json", strings.NewReader(payload))
    if err != nil {
        t.Fatalf("failed to create customer: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("create customer returned %d", resp.StatusCode)
    }

    var created Customer
    if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
        t.Fatalf("failed to decode created customer: %v", err)
    }

    getResp, err := http.Get(fmt.Sprintf("%s/customers/%d", server.URL, created.CustomerID))
    if err != nil {
        t.Fatalf("failed to get created customer: %v", err)
    }
    defer getResp.Body.Close()

    if getResp.StatusCode != http.StatusOK {
        t.Fatalf("get customer returned %d", getResp.StatusCode)
    }

    updatePayload := `{
        "first_name":"Coverage",
        "last_name":"User",
        "email":"coverage.user.updated@example.com",
        "phone":"555-0999",
        "ssn":"111-22-3333",
        "customer_status":"active"
    }`

    req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/customers/%d", server.URL, created.CustomerID), bytes.NewBufferString(updatePayload))
    if err != nil {
        t.Fatalf("failed to build update request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    putResp, err := http.DefaultClient.Do(req)
    if err != nil {
        t.Fatalf("failed to update customer: %v", err)
    }
    defer putResp.Body.Close()

    if putResp.StatusCode != http.StatusOK {
        t.Fatalf("update customer returned %d", putResp.StatusCode)
    }

    body, err := io.ReadAll(putResp.Body)
    if err != nil {
        t.Fatalf("failed to read update response: %v", err)
    }

    var updated Customer
    if err := json.Unmarshal(body, &updated); err != nil {
        t.Fatalf("failed to decode updated customer: %v", err)
    }
    if updated.Email != "coverage.user.updated@example.com" {
        t.Fatalf("update response email mismatch: %s", updated.Email)
    }

    // Cleanup: keep the test database stable for future executions.
    _, err = dbConn.Exec("DELETE FROM customers WHERE customer_id = $1", created.CustomerID)
    if err != nil {
        t.Fatalf("failed to cleanup customer row: %v", err)
    }
}
