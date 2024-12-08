package models

import (
	"sort"

	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
)


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