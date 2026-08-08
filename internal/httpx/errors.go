// Package httpx holds the small HTTP helpers shared by the service entrypoints.
package httpx

import (
    "errors"
    "log"
    "net/http"

    "github.com/lib/pq"
)

// uniqueViolation is the PostgreSQL SQLSTATE for a unique constraint conflict.
const uniqueViolation = "23505"

// IsUniqueViolation reports whether a write failed on a unique constraint, so a
// handler can answer with a conflict status without echoing the driver message.
func IsUniqueViolation(err error) bool {
    var pqErr *pq.Error
    return errors.As(err, &pqErr) && string(pqErr.Code) == uniqueViolation
}

// WriteServerError logs the underlying failure with the operation that produced
// it and returns a generic message so driver text, schema names, and row values
// never reach the client.
func WriteServerError(w http.ResponseWriter, operation string, err error) {
    log.Printf("%s: %v", operation, err)
    http.Error(w, "internal server error", http.StatusInternalServerError)
}
