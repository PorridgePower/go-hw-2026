package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var ch rune
	var cnt int
	var result strings.Builder
	isEscaped := false

	for _, c := range input {
		if isEscaped {
			if unicode.IsDigit(c) || c == '\\' {
				result.WriteRune(ch)
				ch = c
				isEscaped = false
				continue
			}
			return "", ErrInvalidString
		}

		switch {
		case unicode.IsDigit(c) && ch != 0:
			cnt, _ = strconv.Atoi(string(c))
			result.WriteString(strings.Repeat(string(ch), cnt))
			ch = 0
		case unicode.IsDigit(c):
			return "", ErrInvalidString
		case c == '\\':
			isEscaped = true
			continue
		default:
			if ch != 0 {
				result.WriteRune(ch)
			}
			ch = c
		}
	}

	if isEscaped {
		return "", ErrInvalidString
	}

	if ch != 0 {
		result.WriteRune(ch)
	}

	out := result.String()
	return out, nil
}
