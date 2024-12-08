package models

import (
	"fmt"

	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
)

// param  : query (normalized && stemmed)
// return : array os strings (sorted documents names based on the similarity(cosine))
func (data * Data) VSM(query []string)([]string, error) {

	// TFxIDF for each term in the normalized & stemmed query
    vectorQ , err := queryprocess.QueryWeight(query, data.TDM_MATRIX)
	if err != nil {
		return nil , err
	}

    // vector for each document
    results := make(map[string]float64, 0)

    // Calculate cosine similarity.
    // range over all available documents 
    for _, doc := range data.TDM_MATRIX.Documents {

		// create a vector for the document with a length of the unique terms 
        vectorDoc := make([]float64, len(data.TDM_MATRIX.Terms))

        // copy the TFxIDF into the vector for each [term][document]
        for i, term := range data.TDM_MATRIX.Terms {
            vectorDoc[i] = data.TDM_MATRIX.Matrix[term][doc] 
        }

        // Calculate cosine similarity cosine of the document vector and query vector 
        cosineSimilarity, err := utils.CosineSimilarity(vectorQ, vectorDoc)
        if err != nil {
            return nil, fmt.Errorf("failed to calculate cosine similarity: %w", err)
        }

        // Add the document to results if similarity exceeds the threshold
        if cosineSimilarity > 0 {
            results[doc] = cosineSimilarity
        }
        
    }
    // sort the similarity and return a list of an ordered documents as the result 
    finalResult := utils.SortResults(results)
    return finalResult, nil
}