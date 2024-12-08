package models

import (
	"sort"

	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
)


func (data * Data) LSI(query []string) ([]string, error) {

    // calculate the query weight for each term TFxIDF
    processed , err := queryprocess.QueryWeight(query, data.TDM_MATRIX)
	if err != nil {
		return nil , err
	}

    //sort data 
    sort.Strings(data.TDM_MATRIX.Terms)
    sort.Strings(data.TDM_MATRIX.Documents)


    // 2D array of documents with terms instead of map[string]map[string]float64
    TDM := make([][]float64, len(data.TDM_MATRIX.Terms)) 
    for i := range TDM {
        TDM[i] = make([]float64, len(data.TDM_MATRIX.Documents)) 
    }

    // move the TFxIDF to the new 2D matrix 
    for docIdx , doc := range data.TDM_MATRIX.Documents {
        for termIdx , term := range data.TDM_MATRIX.Terms {
            if value, exists := data.TDM_MATRIX.Matrix[term][doc]; exists {
                TDM[termIdx][docIdx] = value 
            } else {
                TDM[termIdx][docIdx] = 0.0 
            }
        }
    }


    // application od SVD that returns U x SIGMA x VT 
    U, SIGMA, VT := utils.SVD(TDM)

    // transform the query to the new space
    newQuery := queryprocess.TransformQueryAlt( processed , U , SIGMA , 5)

    // calculate similarities 
    result := utils.CalculateSimilarities(newQuery , VT , 5)
    documents := make([]string ,0)

    // sort the result 
    sort.Slice(result, func(i, j int) bool {
        return result[i] > result[j]
    })

    // grab the result based on the order of docs and similarity > 0 
    for idx , cos := range result {
        if cos > 0 {
            documents = append(documents , data.TDM_MATRIX.Documents[idx])
        }
    }

    return documents, nil
}