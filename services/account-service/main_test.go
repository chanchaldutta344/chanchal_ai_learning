package main

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "bankcard-learning/internal/db"
)

func TestAccountServiceEndpoints(t *testing.T) {
    dbConn, err := db.InitDB()
    if err != nil {
        t.Fatalf("failed to open database connection: %v", err)
    }
    defer dbConn.Close()

    mux := http.NewServeMux()

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte(`{"status":"ok"}`))
    })

    mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            listAccounts(w, dbConn)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
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
        if suffix == "" {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        split := strings.Split(suffix, "/")
        if len(split) == 1 {
            id := split[0]
            if id == "" {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            accountID := 0
            _, err := w.Write([]byte(""))
            _ = accountID
            _ = err
        }
    })

    server := httptest.NewServer(mux)
    defer server.Close()

    resp, err := http.Get(server.URL + "/health")
    if err != nil {
        t.Fatalf("health endpoint failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("health status expected 200 got %d", resp.StatusCode)
    }

    resp2, err := http.Get(server.URL + "/accounts")
    if err != nil {
        t.Fatalf("account list failed: %v", err)
    }
    defer resp2.Body.Close()

    if resp2.StatusCode != http.StatusOK {
        t.Fatalf("account list expected 200 got %d", resp2.StatusCode)
    }

    var accounts []Account
    if err := json.NewDecoder(resp2.Body).Decode(&accounts); err != nil {
        t.Fatalf("failed to decode accounts response: %v", err)
    }

    if len(accounts) == 0 {
        t.Fatalf("expected seeded accounts from the DDL sample data")
    }
}
