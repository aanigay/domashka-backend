package custom_errors

import "errors"

var (
	ErrAddressNotInRussia = errors.New("geo is not in Russia")
	InvalidAddress        = errors.New("invalid geo")
)
