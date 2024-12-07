package utils

import "sort"


func SortResults(results map[string]float64) []string {
    type docSimilarity struct {
        DocID    string
        Similarity float64
    }

    similaritySlice := make([]docSimilarity, 0, len(results))
    for docID, similarity := range results {
        similaritySlice = append(similaritySlice, docSimilarity{DocID: docID, Similarity: similarity})
    }

    sort.Slice(similaritySlice, func(i, j int) bool {
        return similaritySlice[i].Similarity > similaritySlice[j].Similarity
    })

    sortedDocs := make([]string, len(similaritySlice))
    for i, entry := range similaritySlice {
        sortedDocs[i] = entry.DocID
    }

    return sortedDocs
}