package utils

import (
	"strings"

	"github.com/BigBr41n/scientific-IR/internals/types"
)


func Normalize(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))
	// Remove punctuation and special characters
	word = types.PunctuationRegex.ReplaceAllString(word, "")
	// Remove numbers
	word = types.NumberRegex.ReplaceAllString(word, "")
	
	return word 
}

