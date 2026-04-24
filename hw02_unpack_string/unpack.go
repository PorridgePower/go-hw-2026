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
		if unicode.IsNumber(c) && ch != 0 {
			cnt, _ = strconv.Atoi(string(c))
			// if ch == 0 {
			// 	return "", ErrInvalidString
			// } else {
			result.WriteString(strings.Repeat(string(ch), cnt))
			cnt = 0
			ch = 0
			// }
		} else if unicode.IsNumber(c) {
			return "", ErrInvalidString
		} else {
			result.WriteRune(ch)
			ch = c
		}
	}

	return result.String(), nil
}
