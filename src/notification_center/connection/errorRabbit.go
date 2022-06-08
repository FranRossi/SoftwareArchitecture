package connections

import "errors"

var (
	ErrInvalidQuantity error = errors.New("worker listen quantity must be 0 or positive")
)
