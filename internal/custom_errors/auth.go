package custom_errors

import "fmt"

var (
	ErrUserExists = fmt.Errorf("user already exists")
)
