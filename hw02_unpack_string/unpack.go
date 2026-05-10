package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	fmt.Printf("Unpack: input=%q\n", input)
	if len(input) == 0 {
		fmt.Printf("Unpack: empty input => \"\"\n")
		return "", nil
	}

	var ch rune
	var cnt int
	var result strings.Builder
	var isEscaped bool = false

	for idx, c := range input {
		fmt.Printf("Unpack: idx=%d c=%q ch(prev)=%q cnt=%d isEscaped=%v\n", idx, c, ch, cnt, isEscaped)

		if isEscaped {
			if unicode.IsDigit(c) || c == '\\' {
				fmt.Printf("Unpack: escaped append %q\n", c)
				result.WriteRune(ch)
				// result.WriteRune(c)
				ch = c
				isEscaped = false
				continue
			}
			fmt.Printf("Unpack: invalid escape sequence, got %q\n", c)
			return "", ErrInvalidString
		}

		switch {
		case unicode.IsDigit(c) && ch != 0:
			cnt, _ = strconv.Atoi(string(c))
			fmt.Printf("Unpack: repeat %q %d times\n", ch, cnt)
			result.WriteString(strings.Repeat(string(ch), cnt))
			cnt = 0
			ch = 0
		case unicode.IsDigit(c):
			fmt.Printf("Unpack: digit %q without previous rune => invalid\n", c)
			return "", ErrInvalidString
		case c == '\\':
			isEscaped = true
			fmt.Printf("Unpack: start escape after backslash\n")
			continue
		default:
			if ch != 0 {
				fmt.Printf("Unpack: write pending rune %q\n", ch)
				result.WriteRune(ch)
			}
			ch = c
			fmt.Printf("Unpack: set pending rune to %q\n", ch)
		}
	}

	if ch != 0 {
		fmt.Printf("Unpack: write last pending rune %q\n", ch)
		result.WriteRune(ch)
	}

	out := result.String()
	fmt.Printf("Unpack: output=%q\n", out)
	return out, nil
}
