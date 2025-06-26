package utils

import "unicode"

func BoolToYesNo(val bool) string {
	if val {
		return "Yes"
	}
	return "No"
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
