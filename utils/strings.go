package utils

import "unicode/utf8"

// StringLen returns the length of one string
func StringLen(text string) int {
	return utf8.RuneCountInString(text)
}
