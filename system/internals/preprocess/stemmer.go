package preprocess

import (
	"github.com/reiver/go-porterstemmer"
)

func StemWords(word string) (string) {
	
	stemmedWord := porterstemmer.StemString(word)

	return stemmedWord
}