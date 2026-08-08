package testing

import "os"

// EnvOrDefault returns a runtime override when present, otherwise a safe
// development-time default for local database-oriented service tests.
func EnvOrDefault(name string, fallback string) string {
    if value, ok := os.LookupEnv(name); ok && value != "" {
        return value
    }
    return fallback
}
