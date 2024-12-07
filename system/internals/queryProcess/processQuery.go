package queryprocess

import (
	"math"
	"sort"
	"strings"

	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	"github.com/BigBr41n/scientific-IR/internals/types"
	"github.com/BigBr41n/scientific-IR/internals/utils"
	"github.com/BigBr41n/scientific-IR/internals/weighting"
	"gonum.org/v1/gonum/mat"
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
		word = utils.Normalize(word)
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

		// stem 
		word = preprocess.StemWords(word)
		if _ , exists := qIDF[word]; exists {
			qIDF[word] = Weighting{
				Tf: qIDF[word].Tf + 1,
                Idf: qIDF[word].Idf,
                TFIDF: 0.0,
			}
			continue
		}

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


	for word := range qIDF {
		qIDF[word] = Weighting{
			Tf: qIDF[word].Tf,
            Idf: qIDF[word].Idf,
			TFIDF: qIDF[word].Idf * (1 + math.Log(float64(qIDF[word].Tf))),
		}
	}

	queryResult := make([]float64, 0)
	isMatched := false

	// sort for cosine similarity
	sort.Strings(TDM.Terms)
	sort.Strings(TDM.Documents)

	for _ , term  := range TDM.Terms {
		// matched or not
		isMatched = false
		for word , data := range qIDF {
			if word == term {
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



// LSI 
func TransformQueryAlt(query []float64, U *mat.Dense, Sigma *mat.Dense, k int) *mat.Dense {
	// Step 1: Extract U[t,k]
	rows, _ := U.Dims()
	reducedU := mat.NewDense(rows, k, nil)
	reducedU.Copy(U.Slice(0, rows, 0, k))

	// Step 2: Extract Σ[k,k] and compute Σ[k,k]⁻¹
	sigmaInv := mat.NewDense(k, k, nil)
	for i := 0; i < k; i++ {
		value := Sigma.At(i, i)
		if value != 0 {
			sigmaInv.Set(i, i, 1/value) // Take reciprocal
		} else {
			panic("Singular value is zero, cannot compute inverse")
		}
	}

	// Step 3: Multiply Qold * U[t,k]
	queryVec := mat.NewDense(1, len(query), query) // Query as row vector
	intermediate := mat.NewDense(1, k, nil)
	intermediate.Mul(queryVec, reducedU)

	// Step 4: Multiply the result by Σ⁻¹
	finalQuery := mat.NewDense(1, k, nil)
	finalQuery.Mul(intermediate, sigmaInv)

	return finalQuery
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



