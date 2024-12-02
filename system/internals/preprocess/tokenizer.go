package preprocess

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BigBr41n/scientific-IR/internals/utils"
)

// enriched posting node
type PositionNode struct {
    Position       int16             // position of the word in the Document 
    Next         * PositionNode    // points to the next position node 
    Jump         * PositionNode     // to reduce the comparison 
}

// Posting node for each term (linked list)
type PostingNode struct {
    DocID          string               // string! i'll use the doc name as an ID currently 
    Positions    * PositionNode         // point to all word positions in the document e.g : "term" : ["name of the doc", [1]->[2]->[3]]]->["name of doc2..."]... 
    PositionsTail * PositionNode        // tail of the list of positions
    Jump         * PostingNode          // to reduce the comparison    
    Next         * PostingNode          // points to the next posting node
}   


// The inverted Index type  
type InvertedIndex map[string]*PostingNode 


// Tokenizer Struct 
type Tokenizer struct {
    docsPath string
}

// interface 
type TokenizerI interface {
    ProcessFiles() (InvertedIndex , error)
}



func NewTokenizer(docsPath string) TokenizerI {
    return &Tokenizer{
        docsPath: docsPath,
    }
}



// tokenize and return the local inverted index for each document
func tokenize(entry fs.DirEntry, wg * sync.WaitGroup, indexChan chan <- InvertedIndex, path string ){
    defer wg.Done()

    // local inverted index
    localIndex :=  make(InvertedIndex)

    // Open the file
    filePath := filepath.Join(path, entry.Name())
    file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
    if err != nil {
        log.Print("Error opening file",entry.Name())
        return
    }
    defer file.Close()


    // Load stop words from the file
	stopWords, err := utils.LoadStopWords()
	if err != nil {
		log.Printf("Error loading stop words, the TDM will contain all words")
	}


    // Read file line by line
    wordCount := 0
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // Tokenize the line 
        words := strings.Fields(scanner.Text())

        for _, word := range words {
            wordCount++
  
            // normalize words
            word = utils.Normalize(word)

            // remove stop words
		    if _, isStopWord := stopWords[word]; isStopWord {
			    continue
		    }

            // stem each word 
            word := StemWords(word)

            // position node of the current word
            newPosNode := &PositionNode{
                Position: int16(wordCount),
                Next: nil,
                Jump: nil,
            }

            if currentPosting , exists := localIndex[word]; !exists {
                localIndex[word] = &PostingNode{
                    DocID: entry.Name(),
                    Positions: newPosNode,
                    PositionsTail: newPosNode,
                    Next: nil,
                    Jump: nil,
                }
            } else {
                currentPosting.PositionsTail.Next = newPosNode; 
                currentPosting.PositionsTail = newPosNode;
            }
        }

        // Check if there was a problem reading the file
        if err := scanner.Err(); err != nil {
            log.Printf("Error reading file: %s", entry.Name()) 
            return
        }
    }

    // Send the local index to the main goroutine
    indexChan <- localIndex
}


// open each file in the dir and pass it to Tokenize  function 
// to get the local inverted index for each doc and then merge them in one inverted index 
func (dp * Tokenizer) ProcessFiles() (InvertedIndex , error) {
    // wait group 
    var wg sync.WaitGroup

    // Create a new inverted index (this is the global one that will returned)
    index := make(InvertedIndex)


    // Read Dir 
    entries , err := os.ReadDir(dp.docsPath)
    if err!= nil {
        log.Fatal(err)
    }

    // channel to pass the local indexes from go routines 
    indexChan := make(chan InvertedIndex, len(entries))
    
    // pass each file to a goroutine to handle it 
    // Iterate over each file in the directory
    for _, entry := range entries {
        wg.Add(1)
        if !entry.IsDir() && filepath.Ext(entry.Name()) == ".txt" {
            go tokenize(entry, &wg, indexChan , dp.docsPath)
        }
    }

    // wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(indexChan)
	}()


    // merge inverted indexes into one global inverted index
    for localIndex := range indexChan {
		for term, postingList := range localIndex {
			if _, exists := index[term]; !exists {
				// If term doesn't exist, add it
				index[term] = postingList
			} else {
				// Merge the posting lists
				current := index[term]
				for current.Next != nil {
					current = current.Next
				}
				current.Next = postingList
			}
		}
	}

	return index, nil

}

