package models

import (
	"math"
	"sort"

	"github.com/BigBr41n/scientific-IR/internals/utils"
)

// BM25 model
func (data * Data) BM25(query []string , k int , b float64) ([]string , error) {
    //sort data
    sort.Strings(data.TDM_MATRIX.Terms)
    sort.Strings(data.TDM_MATRIX.Documents)

    // avg doc len
    avgDocLen := float64(len(data.TDM_MATRIX.Documents)) / float64(len(data.TDM_MATRIX.Terms))

    results := make(map[string]float64, len(data.TDM_MATRIX.Documents))

    for _ , document := range data.TDM_MATRIX.Documents {
        documentLen := float64(data.TDM_MATRIX.DocWordCount[document])
        documentScore := 0.0

        for _, term := range query {
            posting := (*data.InvertedIndex)[term]
            tf := 0.0
            for posting != nil && posting.DocID != document {
                posting = posting.Next
            }
            if posting!= nil {
                for pos := posting.Positions; pos!= nil; pos = pos.Next {
                    tf++
                }
            }
            
            numDocs := len(data.TDM_MATRIX.Documents)
            docWithTerm := len(data.TDM_MATRIX.Matrix[term]) 
            idf := math.Log(float64(numDocs-docWithTerm) + 0.5 / float64(docWithTerm)+ 0.5)

            // Compute BM25 score for the term in this document
            documentScore += idf * ((tf * float64(k + 1)) / (tf + float64 (float64(k) * (float64((1 - b)) + b * (documentLen/avgDocLen)))))
        }
        results[document] = documentScore
    }

    sortedDocs := utils.SortResults(results)
    return sortedDocs, nil
    // return results , nil 
}