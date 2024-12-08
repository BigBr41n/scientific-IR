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
    // the number of documents
    numDocs := len(data.TDM_MATRIX.Documents)

    results := make(map[string]float64, len(data.TDM_MATRIX.Documents))

    // range over all documents we have 
    for _ , document := range data.TDM_MATRIX.Documents {
        // store the document length 
        documentLen := float64(data.TDM_MATRIX.DocWordCount[document])

        // initiate the score 
        documentScore := 0.0

        // range over the query terms 
        for _, term := range query {
            // calculating the term frequency in the current document using inverted index 
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
        
            // Calculate the number of documents containing the term
            docWithTerm := len(data.TDM_MATRIX.Matrix[term]) 

            // Calculate IDF of the term
            idf := math.Log(float64(numDocs-docWithTerm) + 0.5 / float64(docWithTerm)+ 0.5)

            // Compute BM25 score for the term in this document using the formula of OKAPI score
            documentScore += idf * ((tf * float64(k + 1)) / (tf + float64 (float64(k) * (float64((1 - b)) + b * (documentLen/avgDocLen)))))
        }
        // store the score for each document with the terms of the query
        results[document] = documentScore
    }

    // sort the scores and get back the list of documents
    sortedDocs := utils.SortResults(results)
    return sortedDocs, nil
}