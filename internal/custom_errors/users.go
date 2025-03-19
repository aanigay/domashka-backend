package custom_errors

import "fmt"

var (
	ErrUserNotFound = fmt.Errorf("user not found")
)
