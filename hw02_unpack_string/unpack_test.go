package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
		{
			input:    "aaaъ3",
			expected: "aaaъъъ",
		},
		{
			input:    "本5",
			expected: "本本本本本",
		},
		{
			input:    " 本5",
			expected: " 本本本本本",
		},
		{
			input:    "  a aa0b ",
			expected: "  a ab ",
		},
		{
			input:    " 9",
			expected: "         ",
		},
		{
			input:    " 9`",
			expected: "         `",
		},
		{
			input:    "@#$%^&*)_",
			expected: "@#$%^&*)_",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {

	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `qwe \\\3`,
			expected: `qwe \3`,
		},
		{
			input:    `qwe\ \\3`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `\\\\\`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `\`,
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    `qwe\456789`,
			expected: "",
			err:      ErrInvalidString,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
