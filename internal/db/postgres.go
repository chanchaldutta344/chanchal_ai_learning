package db

import (
    "database/sql"
    "fmt"
    "os"

    _ "github.com/lib/pq"
)

// InitDB centralizes the PostgreSQL configuration for all service entrypoints.
// It reads runtime environment variables, builds a DSN, and verifies the server
// can accept a connection before returning a database handle.
func InitDB() (*sql.DB, error) {
    host := os.Getenv("PGHOST")
    if host == "" {
        host = "localhost"
    }
    port := os.Getenv("PGPORT")
    if port == "" {
        port = "5432"
    }
    user := os.Getenv("PGUSER")
    if user == "" {
        user = "appuser"
    }
    password := os.Getenv("PGPASSWORD")
    if password == "" {
        password = "appuser123"
    }
    dbname := os.Getenv("PGDATABASE")
    if dbname == "" {
        dbname = "appdb"
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
