package preprocess

import (
	"log"
	"sync"

	"github.com/jdkato/prose/v2"
	"github.com/reiver/go-porterstemmer"
)

var lemmaCache = sync.Map{}

func LemmatizeStemWords(word string , stemORstemLem int) string {
    //log.Println(word)

	if stemORstemLem == 0 {
		// Apply Porter stemming only
		return porterstemmer.StemString(word)
	}

    // Check if the word is cached
    if cached, ok := lemmaCache.Load(word); ok {
        return cached.(string)
    }

    // Create a prose document
    doc, err := prose.NewDocument(word)
    if err != nil {
        log.Fatalf("Error creating document for word '%s': %v", word, err)
    }

    // Get tokens from the document
    tokens := doc.Tokens()
    if len(tokens) == 0 {
        //log.Printf("No tokens generated for word '%s'. Returning original word.", word)
        return word // Fallback to the original word
    }

    // Use the first token's text as the lemma
    lemma := tokens[0].Text

    // Cache the lemma
    lemmaCache.Store(word, lemma)

    // Apply Porter stemming
    stemmedWord := porterstemmer.StemString(lemma)

    return stemmedWord
}