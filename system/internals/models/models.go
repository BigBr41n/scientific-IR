package models

import (
	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	"github.com/BigBr41n/scientific-IR/internals/types"
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

