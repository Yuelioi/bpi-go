// Package bpierr contains shared error implementations used across the root
// and domain packages. Public callers use the aliases exposed by package bpi.
package bpierr

import "fmt"

// ParameterError describes an invalid caller-supplied value.
type ParameterError struct {
	Field   string
	Message string
}

func (e *ParameterError) Error() string {
	return fmt.Sprintf("bpi: invalid parameter %q: %s", e.Field, e.Message)
}
