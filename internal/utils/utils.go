package utils

import (
	"strconv"
	"unicode"
)

func ValidLuhn(number string) bool {
	sum := 0
	nDigits := len(number)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		r := rune(number[i])
		if !unicode.IsDigit(r) {
			return false
		}
		digit, _ := strconv.Atoi(string(r))
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
