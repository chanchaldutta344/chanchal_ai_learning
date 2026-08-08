// Package money centralizes the exact decimal representation used for the
// DECIMAL(12, 2) monetary columns in the schema.
package money

import "github.com/shopspring/decimal"

// Amount is the exact decimal type used for every monetary field so credit
// limits, balances, and payment amounts never drift through binary floats.
type Amount = decimal.Decimal

func init() {
    // Monetary values stay JSON numbers instead of quoted strings.
    decimal.MarshalJSONWithoutQuotes = true
}

// IsPositive reports whether an amount is strictly greater than zero.
func IsPositive(a Amount) bool {
    return a.IsPositive()
}
