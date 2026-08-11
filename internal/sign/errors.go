// Package sign implements deterministic signing primitives used by Bilibili
// request policies.
package sign

import "fmt"

// InvalidError describes invalid input to a signing primitive.
type InvalidError struct {
	Field   string
	Message string
}

func (e *InvalidError) Error() string {
	return fmt.Sprintf("sign: invalid %s: %s", e.Field, e.Message)
}
