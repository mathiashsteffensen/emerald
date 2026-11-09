package core

import (
	"unicode"
	"unicode/utf8"
)

// IsConstantName reports whether name is a Ruby constant name: its first
// character is uppercase (including Other_Uppercase such as Ⅰ) or titlecase,
// ASCII or not.
func IsConstantName(name string) bool {
	firstChar, _ := utf8.DecodeRuneInString(name)

	return unicode.IsUpper(firstChar) ||
		unicode.IsTitle(firstChar) ||
		unicode.Is(unicode.Other_Uppercase, firstChar)
}
