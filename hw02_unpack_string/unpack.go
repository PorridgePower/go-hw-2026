package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if len(input) == 0 {
		return "", nil
	}
	var ch rune
	var cnt int
	var result []rune
	for _, c := range input {
		if unicode.IsNumber(c) && ch != 0 {
			cnt, _ = strconv.Atoi(string(c))
			if cnt == 0 {
				result = result[:len(result)-1]
			} else {
				for range cnt - 1 {
					result = append(result, ch)
				}
			}
			cnt = 0
			ch = 0
		} else if unicode.IsNumber(c) {
			return "", ErrInvalidString
		} else {
			result = append(result, c)
			ch = c
		}
	}

	return string(result), nil
}
