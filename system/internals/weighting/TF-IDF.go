package weighting

import (
	"math"

	"github.com/BigBr41n/scientific-IR/internals/types"
)

// TF= Term Frequency in Document / Total Terms in Document
// IDF= log(Total Documents / Number of Documents Containing the Term + 1 )



func CalculateTFIDF(tdm * types.TDM) {
	totalDocs := len(tdm.Documents)

	// Calculate TF-IDF for each term-document pair
	for term, docFreqMap := range tdm.Matrix {
		// Calculate IDF
		docCountWithTerm := len(docFreqMap)
		idf := math.Log(float64(totalDocs) / float64(docCountWithTerm)+1)

		// Update TF-IDF weights
		for doc, termFrequency := range docFreqMap {
			tf := termFrequency / float64(tdm.DocWordCount[doc]) 
			tdm.Matrix[term][doc] = tf * idf
		}
	}
}


func CalculateQueryIDF(word string, tdm *types.TDM) float64 {
	totalDocs := len(tdm.Documents)

	// check if the word exists in the matrix
	docFreqMap, exists := tdm.Matrix[word]
	if !exists || len(docFreqMap) == 0 {
		return 0.0
	}


	docCountWithTerm := len(docFreqMap)

	// calculate IDF using the formula: log(totalDocs / (docCountWithTerm))
	idf := math.Log(float64(totalDocs) / float64(docCountWithTerm)+1)
	return idf
}
