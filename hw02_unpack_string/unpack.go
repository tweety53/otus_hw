package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {

	if str == "" {
		return "", nil
	}

	hasInvalidSeq, _ := regexp.MatchString(`^[0-9].|\\[^(0-9|\\)]|[^\\]\d{2,}|\\(\\{2})$|^\\$`, str)

	if hasInvalidSeq {
		return "", ErrInvalidString
	}

	var (
		b       strings.Builder
		prev    rune
		escaped bool
	)

	for _, r := range str {
		if r == '\\' && !escaped {
			escaped = true
			continue
		}

		if (!unicode.IsDigit(r) && !unicode.IsDigit(prev) && prev != 0) || escaped {
			b.WriteRune(prev)
			prev = r
			escaped = false
			continue
		}

		if unicode.IsDigit(r) && !escaped {
			b.WriteString(strings.Repeat(string(prev), int(r-'0')))
			prev = rune(0)
			continue
		}

		prev = r
	}

	if prev != rune(0) {
		b.WriteRune(prev)
	}

	return b.String(), nil
}
