package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BigBr41n/scientific-IR/internals/models"
	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
)

//func printInvertedIndex(index preprocess.InvertedIndex) {
//    for term, postingList := range index {
//        log.Printf("Term: %s", term)
//
//        currentPosting := postingList
//        for currentPosting != nil {
//            log.Printf("  Document: %s", currentPosting.DocID)
//
//            // Traverse the positions list
//            positions := []int16{}
//            currentPosition := currentPosting.Positions
//            for currentPosition != nil {
//                positions = append(positions, currentPosition.Position)
//                currentPosition = currentPosition.Next
//            }
//
//            log.Printf("    Positions: %v", positions)
//
//            // Move to the next posting in the list
//            currentPosting = currentPosting.Next
//        }
//    }
//}



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

	//printInvertedIndex(invertedIndex)
    // create term document matrix 
    TDM:= preprocess.BuildTDM(invertedIndex)
    //log.Println("Term Document Matrix:")
    //preprocess.PrintTDM(TDM)

    stopWords, _ := utils.LoadStopWords();

    query := ""
    fmt.Print("Enter a query: ")
    scanner := bufio.NewScanner(os.Stdin)

	// Scan the next line.
	if scanner.Scan() {
		query = scanner.Text()
	}

    models := models.NewInfoRetrievalModel(TDM , &stopWords, &invertedIndex)

    result, _ := models.ClassicBoolean(query)

    fmt.Printf("the result with intersection : %v\n", result)



	fmt.Print("Enter a query: ")
	// Scan the next line.
	if scanner.Scan() {
		query = scanner.Text()
	}

	result , _ = models.VSM(query)

	fmt.Printf("the result with VSM : %v\n", result)



	fmt.Print("Enter a query: ")
	if scanner.Scan() {
		query = scanner.Text()
	}

	result , _ = models.LSI(query)

	fmt.Printf("the result with LSI : %v\n", result)
}