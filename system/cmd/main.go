package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BigBr41n/scientific-IR/internals/preprocess"
)



func printInvertedIndex(index preprocess.InvertedIndex) {
    for term, postingList := range index {
        log.Printf("Term: %s", term)

        currentPosting := postingList
        for currentPosting != nil {
            log.Printf("  Document: %s", currentPosting.DocID)

            // Traverse the positions list
            positions := []int16{}
            currentPosition := currentPosting.Positions
            for currentPosition != nil {
                positions = append(positions, currentPosition.Position)
                currentPosition = currentPosition.Next
            }

            log.Printf("    Positions: %v", positions)

            // Move to the next posting in the list
            currentPosting = currentPosting.Next
        }
    }
}



func main() {
	log.Println("SCIENTIFIC INFORMATION RETRIEVAL")

	cwd , err := os.Getwd()
	if err!= nil {
        log.Fatal(err)
    }
	log.Println("Current working directory:", cwd)
	dirPath := filepath.Join(cwd,"/data/documents")
	tokenizer := preprocess.NewTokenizer(dirPath)

	invertedIndex , _ := tokenizer.ProcessFiles()

	printInvertedIndex(invertedIndex)
}