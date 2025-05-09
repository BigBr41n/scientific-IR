package preprocess

import (
	"log"

	"github.com/BigBr41n/scientific-IR/internals/types"
	"github.com/BigBr41n/scientific-IR/internals/weighting"
)

// using preprocessed inverted index is much faster than
// reprocess files and extract the number of occurrence for each word
// BuildTDM constructs a Term-Document Matrix (frequency-based) from the InvertedIndex
func BuildTDM(index InvertedIndex) * types.TDM {
	tdm := &types.TDM{
		Matrix:      make(map[string]map[string]float64),
		Terms:       []string{},
		Documents:   []string{},
		DocWordCount: make(map[string]int),
	}

	// Collect all terms and documents
	docSet := make(map[string]struct{}) 
	for term, postingList := range index {

		tdm.Terms = append(tdm.Terms, term)

		currentPosting := postingList
		for currentPosting != nil {
			docID := currentPosting.DocID

			if _, exists := docSet[docID]; !exists {
				docSet[docID] = struct{}{}
			}

			// initialize the term-document map
			if _, exists := tdm.Matrix[term]; !exists {
				tdm.Matrix[term] = make(map[string]float64)
			}

			// count the frequency of the term in the document
			positionNode := currentPosting.Positions
			frequency := 0
			for positionNode != nil {
				frequency++
				positionNode = positionNode.Next
			}

			// update the raw frequency in the TDM
			tdm.Matrix[term][docID] = float64(frequency)

			// Increment total word count for the document
			tdm.DocWordCount[docID] += frequency

			currentPosting = currentPosting.Next
		}
	}

	// Convert the document set to a slice
	for doc := range docSet {
		tdm.Documents = append(tdm.Documents, doc)
	}


	// call TF IDF calculator
	weighting.CalculateTFIDF(tdm)

	return tdm
}


// PrintTDM displays the Term-Document Matrix with TF-IDF weights
func PrintTDM(tdm * types.TDM) {
	log.Printf("%-15s", "Term/Document")
	for _, doc := range tdm.Documents {
		log.Printf("%-15s", doc)
	}
	log.Println()

	for _, term := range tdm.Terms {
		log.Printf("%-15s", term)
		for _, doc := range tdm.Documents {
			log.Printf("%-15.4f", tdm.Matrix[term][doc]) // Print TF-IDF with 4 decimal places
		}
		log.Println()
	}
}