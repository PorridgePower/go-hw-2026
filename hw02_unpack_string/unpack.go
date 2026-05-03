package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	var ch rune
	var cnt int
	var result strings.Builder
	for _, c := range input {
		if unicode.IsDigit(c) && ch != 0 {
			cnt, _ = strconv.Atoi(string(c))
			result.WriteString(strings.Repeat(string(ch), cnt))
			cnt = 0
			ch = 0
		} else if unicode.IsDigit(c) {
			return "", ErrInvalidString
		} else {
			if ch != 0 {
				result.WriteRune(ch)
			}
			ch = c
		}
	}
	if ch != 0 {
		result.WriteRune(ch)
	}

	return result.String(), nil
}
