package custom_errors

import "fmt"

var (
	ErrCartNotFound = fmt.Errorf("cart not found")
)
