package models

import (
	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/types"
)


type Data struct {
	TDM_MATRIX * types.TDM
	StopWords  *  map[string]struct{} 
	InvertedIndex * preprocess.InvertedIndex
}


type IrModels interface {
	ClassicBoolean(query string) ([]string, error)
	//VSM(query string) ([]string , error)
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

//func (data * Data) VSM(query string)([]string, error) {
//
//}
//
//func (data * Data) LSI(query string) ([]string, error) {
//
//}