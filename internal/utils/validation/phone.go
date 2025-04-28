package validation

import "regexp"

func ValidatePhoneNumber(phoneNumber string) bool {
	pattern := `^(\+7|8)?[\s\-\(]*\d{3}[\s\-\)]*\d{3}[\s\-]*\d{2}[\s\-]*\d{2}$`
	r := regexp.MustCompile(pattern)
	return r.MatchString(phoneNumber)
}
