package models

import (
	"fmt"

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
	ClassicBoolean(query string) ([]string, error)
	VSM(query string) ([]string , error)
	//LSI(query string) ([]string , error)
}



func NewInfoRetrievalModel(TDM_MATRIX * types.TDM, StopWords  *  map[string]struct{} , InvertedIndex * preprocess.InvertedIndex  ) IrModels {
	return &Data{
		TDM_MATRIX : TDM_MATRIX ,
		StopWords : StopWords ,
		InvertedIndex : InvertedIndex,
	}
}


func (data * Data) ClassicBoolean(query string) ([]string, error) {
	processed , err := queryprocess.Classic(query, data.StopWords)
	if err != nil {
		return nil , err
	}

    // store the set of matching DocIDs for each term.
    var intersectedDocIDs map[string]struct{}
    isFirstTerm := true // handle the first term differently.

	for _, word := range processed {
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

func (data * Data) VSM(query string)([]string, error) {
    processed , err := queryprocess.QueryWeight(query, data.TDM_MATRIX, data.StopWords)
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
    finalResult := utils.SortResultsBySimilarity(results)

    return finalResult, nil
}



//
//func (data * Data) LSI(query string) ([]string, error) {
//
//}