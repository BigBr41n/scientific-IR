package models

import (
	"fmt"

	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
)


func (data * Data) VSM(query []string)([]string, error) {
    processed , err := queryprocess.QueryWeight(query, data.TDM_MATRIX)
	if err != nil {
		return nil , err
	}

    //log.Println("the len of query : ", len(processed))
    //log.Println("the len of terms : ", len(data.TDM_MATRIX.Terms))
    //log.Println("processed query : ", processed)


    // vector for each document
    results := make(map[string]float64, 0)

    // Calculate cosine similarity.
    for _, doc := range data.TDM_MATRIX.Documents {
        vectorDoc := make([]float64, len(data.TDM_MATRIX.Terms))
        for i, term := range data.TDM_MATRIX.Terms {
            vectorDoc[i] = data.TDM_MATRIX.Matrix[term][doc] 
        }
        //log.Println("processed doc : ", vectorDoc )
        // Calculate cosine similarity
        cosineSimilarity, err := utils.CosineSimilarity(processed, vectorDoc)
        if err != nil {
            return nil, fmt.Errorf("failed to calculate cosine similarity: %w", err)
        }
        // Add the document to results if similarity exceeds the threshold
        if cosineSimilarity > 0 {
            results[doc] = cosineSimilarity
        }
        
    }
    finalResult := utils.SortResults(results)

    return finalResult, nil
}