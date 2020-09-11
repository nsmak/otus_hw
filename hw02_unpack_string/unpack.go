package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}
	var prevRune rune
	var builder strings.Builder
	var isEscapeOn bool

	for _, r := range s {
		if prevRune == 0 {
			if unicode.IsDigit(r) {
				return "", ErrInvalidString
			}
			prevRune = r
			continue
		}
		if prevRune == '\\' && !isEscapeOn {
			isEscapeOn = true
			prevRune = r
			continue
		}
		isEscapeOn = false
		if !unicode.IsDigit(r) {
			builder.WriteRune(prevRune)
			prevRune = r
			continue
		}
		count, err := strconv.Atoi(string(r))
		if err != nil {
			return "", err
		}
		builder.WriteString(strings.Repeat(string(prevRune), count))
		prevRune = 0
	}
	if prevRune != 0 {
		builder.WriteRune(prevRune)
	}

	final := builder.String()
	return final, nil
}
