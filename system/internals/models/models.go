package models

import (
	"fmt"
	"math"
	"sort"

	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/types"
	"github.com/BigBr41n/scientific-IR/internals/utils"
)


type Data struct {
	TDM_MATRIX * types.TDM
	StopWords  *  map[string]struct{} 
	InvertedIndex * preprocess.InvertedIndex
}


type IrModels interface {
	ClassicBoolean(query []string) ([]string, error)
	VSM(query []string) ([]string , error)
	LSI(query []string) ([]string , error)
    BM25(query []string, k int , b float64) ([]string , error)
}



func NewInfoRetrievalModel(TDM_MATRIX * types.TDM, StopWords  *  map[string]struct{} , InvertedIndex * preprocess.InvertedIndex  ) IrModels {
	return &Data{
		TDM_MATRIX : TDM_MATRIX ,
		StopWords : StopWords ,
		InvertedIndex : InvertedIndex,
	}
}


func (data * Data) ClassicBoolean(query []string) ([]string, error) {

    // store the set of matching DocIDs for each term.
    var intersectedDocIDs map[string]struct{}
    isFirstTerm := true // handle the first term differently.

	for _, word := range query {
        posting := (*data.InvertedIndex)[word]

		currentDocIDs := make(map[string]struct{})

        for posting != nil {
            currentDocIDs[posting.DocID] = struct{}{}
            posting = posting.Next
        }

        // intersection with the existing results.
        if isFirstTerm {
            intersectedDocIDs = currentDocIDs
            isFirstTerm = false
        } else {
            for docID := range intersectedDocIDs {
                if _, exists := currentDocIDs[docID]; !exists {
                    delete(intersectedDocIDs, docID)
                }
            }
        }
    }

    // Convert the resulting intersection set to a slice of DocIDs.
    result := make([]string, 0, len(intersectedDocIDs))
    for docID := range intersectedDocIDs {
        result = append(result, docID)
    }

    return result, nil
}



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



func (data * Data) LSI(query []string) ([]string, error) {
    processed , err := queryprocess.QueryWeight(query, data.TDM_MATRIX)
	if err != nil {
		return nil , err
	}

    //sort data 
    sort.Strings(data.TDM_MATRIX.Terms)
    sort.Strings(data.TDM_MATRIX.Documents)


    TDM := make([][]float64, len(data.TDM_MATRIX.Terms)) 
    for i := range TDM {
        TDM[i] = make([]float64, len(data.TDM_MATRIX.Documents)) 
    }

    // move the TFxIDF to the matrix 
    for docIdx , doc := range data.TDM_MATRIX.Documents {
        for termIdx , term := range data.TDM_MATRIX.Terms {
            if value, exists := data.TDM_MATRIX.Matrix[term][doc]; exists {
                TDM[termIdx][docIdx] = value 
            } else {
                TDM[termIdx][docIdx] = 0.0 
            }
        }
    }

    U, SIGMA, VT := utils.SVD(TDM)
    newQuery := queryprocess.TransformQueryAlt( processed , U , SIGMA , 3)


    // calculate similarities 
    result := utils.CalculateSimilarities(newQuery , VT , 3)
    documents := make([]string ,0)


    //log.Println("the length of result : ",len(result)) // should be equal the number of docs

    // sort the result 
    sort.Slice(result, func(i, j int) bool {
        return result[i] > result[j]
    })

    for idx , cos := range result {
        if cos > 0 {
            documents = append(documents , data.TDM_MATRIX.Documents[idx])
        }
    }

    return documents, nil
}



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