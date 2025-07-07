package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

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

func EmailToUsername(email string) string {
	parts := strings.Split(email, "@")
	return parts[0]
}

func StripFields(data interface{}, keysToRemove ...string) []map[string]interface{} {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal error:", err)
		return nil
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		fmt.Println("unmarshal error:", err)
		return nil
	}

	keySet := make(map[string]struct{})
	for _, key := range keysToRemove {
		keySet[key] = struct{}{}
	}

	for i := range result {
		for key := range keySet {
			delete(result[i], key)
		}
	}

	return result
}
