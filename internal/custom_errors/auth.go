package custom_errors

import "fmt"

var (
	ErrUserExists              = fmt.Errorf("user already exists")
	ErrUserIsSpam              = fmt.Errorf("user is spam")
	ErrExpiredTTL              = fmt.Errorf("expired ttl")
	ErrConfirmationNotReceived = fmt.Errorf("confirmation not received")
	ErrPhoneNumberMismatch     = fmt.Errorf("phone number mismatch")
)
