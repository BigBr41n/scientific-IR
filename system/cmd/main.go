package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BigBr41n/scientific-IR/internals/models"
	"github.com/BigBr41n/scientific-IR/internals/preprocess"
	queryprocess "github.com/BigBr41n/scientific-IR/internals/queryProcess"
	"github.com/BigBr41n/scientific-IR/internals/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app        *tview.Application
	mainFlex   *tview.Flex
	queryInput *tview.InputField
	results    *tview.TextView
)

func main() {
	app = tview.NewApplication()

	// Set default black background
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.TitleColor = tcell.ColorWhite
	tview.Styles.BorderColor = tcell.ColorWhite

	// Initialize IR components
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dirPath := filepath.Join(cwd, "/data/documents")
	tokenizer := preprocess.NewTokenizer(dirPath)

	invertedIndex, _ := tokenizer.ProcessFiles()
	TDM := preprocess.BuildTDM(invertedIndex)
	stopWords, _ := utils.LoadStopWords()
	models := models.NewInfoRetrievalModel(TDM, &stopWords, &invertedIndex)

	// Create header
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("SCIENTIFIC INFORMATION RETRIEVAL SYSTEM").
		SetDynamicColors(true).
		SetBackgroundColor(tcell.NewRGBColor(0, 0, 0)) 

	// Create query input
	queryInput = tview.NewInputField().
		SetLabel("[yellow]Enter query >[-]").
		SetFieldWidth(60)

	// Create results view
	results = tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)

	results.SetBackgroundColor(tcell.NewRGBColor(0, 0, 0)) 

	queryInput.SetBackgroundColor(tcell.NewRGBColor(0, 0, 0)) 
	queryInput.SetFieldBackgroundColor(tcell.ColorWhite)   
	queryInput.SetFieldTextColor(tcell.ColorRed)    

	// Handle query submission
	queryInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			query := queryInput.GetText()
			results.Clear()
	
			preprocessedQuery, err := queryprocess.ProcessQuery(query, &stopWords)
			if err != nil {
				fmt.Fprintf(results, "[red]Error processing query: %v\n[-]", err)
				app.SetFocus(queryInput) // Refocus on the input field
				return
			}
	
			// Add a title for the query results
			fmt.Fprintf(results, "[yellow]--- Results for Query: '%s' ---[-]\n\n", query)
	
			// Classic Boolean
			result, _ := models.ClassicBoolean(preprocessedQuery)
			fmt.Fprintf(results, "[cyan]Classic Boolean Results:[-]\n")
			if len(result) > 0 {
				for _, doc := range result {
					fmt.Fprintf(results, "  - %v\n", doc)
				}
			} else {
				fmt.Fprintf(results, "  [gray]No results found.[-]\n")
			}
			fmt.Fprintf(results, "\n")
	
			// VSM
			result, _ = models.VSM(preprocessedQuery)
			fmt.Fprintf(results, "[cyan]Vector Space Model Results:[-]\n")
			if len(result) > 0 {
				for _, doc := range result {
					fmt.Fprintf(results, "  - %v\n", doc)
				}
			} else {
				fmt.Fprintf(results, "  [gray]No results found.[-]\n")
			}
			fmt.Fprintf(results, "\n")
	
			// LSI
			result, _ = models.LSI(preprocessedQuery)
			fmt.Fprintf(results, "[cyan]Latent Semantic Indexing Results:[-]\n")
			if len(result) > 0 {
				for _, doc := range result {
					fmt.Fprintf(results, "  - %v\n", doc)
				}
			} else {
				fmt.Fprintf(results, "  [gray]No results found.[-]\n")
			}
			fmt.Fprintf(results, "\n")
	
			// BM25
			fmt.Fprintf(results, "[cyan]BM25 Results:[-]\n")
			BM25, err := models.BM25(preprocessedQuery, 1, 0.75)
			if err != nil {
				fmt.Fprintf(results, "[red]Error in BM25: %v[-]\n", err)
			} else if len(BM25) > 0 {
				for idx , res := range BM25 {
					if idx < 10 {
						fmt.Fprintf(results, "  - %v\n", res)
					}else {
						break
					}
				}
			} else {
				fmt.Fprintf(results, "  [gray]No results found.[-]\n")
			}
	
			// Add a footer to separate results from input
			fmt.Fprintf(results, "\n[yellow]--- End of Results ---[-]\n\n")
	
			app.SetFocus(queryInput) // Refocus on the input field
		}
	})

	// Layout
	mainFlex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(header, 3, 0, false).
		AddItem(queryInput, 3, 0, true).
		AddItem(results, 0, 1, false)

	if err := app.SetRoot(mainFlex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
