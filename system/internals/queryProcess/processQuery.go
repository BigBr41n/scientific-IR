package queryprocess

import (
	"strings"

	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	"github.com/BigBr41n/scientific-IR/internals/types"
	"github.com/BigBr41n/scientific-IR/internals/utils"
	"github.com/BigBr41n/scientific-IR/internals/weighting"
)



type Weighting struct {
	Tf 		int
	Idf 	float64
	TFIDF   float64
}

type TermTFIDF struct {
    Term   string
    TFIDF  float64
}


// exact matching
func Classic(query string, stopWords * map[string]struct{}) ([]string , error){
	var results []string

	// Load stop words
	// stopWords, err := utils.LoadStopWords()
    // if err != nil {
    //     return nil, err
    // }
	// extract words 
    words := strings.Fields(query)
	// remove stop words 
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(strings.Trim(word, ".,!?\"'")))
        if _, exists := (*stopWords)[word]; exists {
            continue
        }
		// stem 
		word = preprocess.StemWords(word)
        results = append(results, word)
	}

	return results, nil
} 



func QueryWeight(query string, TDM * types.TDM,stopWords * map[string]struct{} )([]float64 , error){
	qIDF := make(map[string]Weighting)
	// Load stop words
	// stopWords, err := utils.LoadStopWords()
    // if err!= nil {
    //     return nil, err
    // }

	// extract words 
	words := strings.Fields(query)

	// remove stop words 
	for _, word := range words {
		
		word = utils.Normalize(word)

		if _, exists := (*stopWords)[word]; exists {
			continue
		}
		if _ , exists := qIDF[word]; exists {
			qIDF[word] = Weighting{
				Tf: qIDF[word].Tf + 1,
                Idf: qIDF[word].Idf,
                TFIDF: 0.0,
			}
			continue
		}

		// stem 
		word = preprocess.StemWords(word)

		IDF := weighting.CalculateQueryIDF(word, TDM)
		if IDF == 0.0 {
			// don't store the term because it don't exist in the unique terms 
			continue
        }

		qIDF[word] = Weighting{
			Tf: 1,
			Idf: IDF,
			TFIDF: 0.0,
		}
	}


	queryLenght := len(qIDF)

	for word := range qIDF {
		qIDF[word] = Weighting{
			Tf: qIDF[word].Tf,
            Idf: qIDF[word].Idf,
			TFIDF: qIDF[word].Idf * (float64(qIDF[word].Tf) / float64(queryLenght)),
		}
	}

	queryResult := make([]float64, 0)
	isMatched := false

	for key := range TDM.Matrix {
		// matched or not
		isMatched = false
		for word , data := range qIDF {
			if word == key {
				queryResult = append(queryResult , data.TFIDF)
				isMatched = true
				break 
			}
		}
		if !isMatched {
			queryResult = append(queryResult, 0.0)
		}
	}

	return queryResult, nil
}



/* func VSMBoolean(query string, TDM * types.TDM)(map[string]float64 , error){
	qIDF := make(map[string]Weighting)
	// Load stop words
	stopWords, err := utils.LoadStopWords()
    if err!= nil {
        return nil, err
    }

	// extract words 
	words := strings.Fields(query)

	// remove stop words 
	for _, word := range words {
		
		word = utils.Normalize(word)

		if _, exists := stopWords[word]; exists && word != "AND" && word != "OR" && word != "NOT"{
			continue
		}
		if _ , exists := qIDF[word]; exists {
			qIDF[word] = Weighting{
				Tf: qIDF[word].Tf + 1,
                Idf: qIDF[word].Idf,
                TFIDF: 0.0,
			}
			continue
		}

		// stem 
		word = preprocess.StemWords(word)

		IDF := weighting.CalculateQueryIDF(word, TDM)
		if IDF == 0.0 && word != "AND" && word != "OR" && word != "NOT" {
			// don't store the term because it don't exist in the unique terms 
			continue
        }

		qIDF[word] = Weighting{
			Tf: 1,
			Idf: IDF,
			TFIDF: 0.0,
		}
	}


	countQueryLen := 0 
	for key := range qIDF {
		if key == "AND" || key == "OR" || key == "NOT" {
			continue
		}
		countQueryLen++
	}

	if (countQueryLen == 0) {
		return map[string]float64{}, nil
	}

	queryLenght := len(qIDF)


	for word := range qIDF {
		if word == "AND" || word == "OR" || word == "NOT" {
            qIDF[word] = Weighting{
				Tf: qIDF[word].Tf,
				Idf: qIDF[word].Idf,
				TFIDF: 0.0,
			}
        }
		qIDF[word] = Weighting{
			Tf: qIDF[word].Tf,
            Idf: qIDF[word].Idf,
			TFIDF: qIDF[word].Idf * (float64(qIDF[word].Tf) / float64(queryLenght)),
		}
	}

	queryResult := make(map[string]float64, 0)

	for word , data := range qIDF {
		queryResult[word] =  data.TFIDF
	}

	return queryResult, nil
} */



