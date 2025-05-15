package types

import (
	"fmt"
	"regexp"
	"strings"
)

var digitsOnly = regexp.MustCompile(`\d`)

// FormatPhoneNumber преобразует указатель на любой формат телефона вида
// +7(XXX) XXX XX-XX
// Если вход nil или в нём меньше 10 цифр — возвращает исходный указатель.
func FormatPhoneNumber(input *string) *string {
	if input == nil {
		return nil
	}
	raw := *input

	// Извлекаем все цифры
	parts := digitsOnly.FindAllString(raw, -1)
	if len(parts) < 10 {
		// Недостаточно цифр — возвращаем без изменений
		return input
	}

	// Берём последние 10 цифр
	digits := parts[len(parts)-10:]
	area := strings.Join(digits[0:3], "") // XXX
	mid := strings.Join(digits[3:6], "")  // XXX
	p1 := strings.Join(digits[6:8], "")   // XX
	p2 := strings.Join(digits[8:10], "")  // XX

	// Формируем ссылку и отображаемый текст
	formatted := fmt.Sprintf("+7(%s) %s %s-%s", area, mid, p1, p2)

	return &formatted
}
