package weighting

import (
	"math"

	"github.com/BigBr41n/scientific-IR/internals/preprocess"
)

// TF= Term Frequency in Document / Total Terms in Document
// IDF= log(Total Documents / Number of Documents Containing the Term)



func CalculateTFIDF(tdm * preprocess.TDM) {
	totalDocs := len(tdm.Documents)

	// Calculate TF-IDF for each term-document pair
	for term, docFreqMap := range tdm.Matrix {
		// Calculate IDF
		docCountWithTerm := len(docFreqMap)
		idf := math.Log(float64(totalDocs) / float64(docCountWithTerm))

		// Update TF-IDF weights
		for doc, termFrequency := range docFreqMap {
			tf := termFrequency / float64(tdm.DocWordCount[doc]) 
			tdm.Matrix[term][doc] = tf * idf
		}
	}
}