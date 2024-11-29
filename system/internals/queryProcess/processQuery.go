package queryprocess

import (
	"strings"

	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
	"github.com/BigBr41n/scientific-IR/internals/weighting"
)

// exact matching
func Classic(query string) ([]string , error){
	var results []string

	// Load stop words
	stopWords, err := utils.LoadStopWords()
    if err != nil {
        return nil, err
    }
	// extract words 
    words := strings.Fields(query)
	// remove stop words 
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(strings.Trim(word, ".,!?\"'")))
        if _, exists := stopWords[word]; exists {
            continue
        }
		// stem 
		word = preprocess.StemWords(word)
        results = append(results, word)
	}

	return results, nil
} 



func QueryWeight(query string, TDM * preprocess.TDM)(map[string]float64 , error){
	qIDF := make(map[string]float64)
	// Load stop words
	stopWords, err := utils.LoadStopWords()
    if err!= nil {
        return nil, err
    }

	// extract words 
	words := strings.Fields(query)

	// remove stop words 
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(strings.Trim(word, ".,!?\"'")))
		if _, exists := stopWords[word]; exists {
			continue
		}
		if _ , exists := qIDF[word]; exists {
			qIDF[word] = qIDF[word] + 1
			continue
		}
		// stem 
		word = preprocess.StemWords(word)
		qIDF[word] = 1
	}

	for word , value := range qIDF {
		IDF := weighting.CalculateQueryIDF(word, TDM)
		if IDF == 0.0 {
            continue
        }
		qIDF[word] = IDF * value
	}

	return qIDF, nil
}




