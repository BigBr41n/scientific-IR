package preprocess

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// enriched posting node
type PositionNode struct {
    position       int16             // position of the word in the Document 
    next         * PositionNode    // points to the next position node 
    jump         * PositionNode     // to reduce the comparison 
}

// Posting node for each term (linked list)
type PostingNode struct {
    docId          string               // string! i'll use the doc name as an ID currently 
    positions    * PositionNode     // point to all word positions in the document e.g : "term" : ["name of the doc", [1]->[2]->[3]]]->["name of doc2..."]... 
    jump         * PostingNode     // to reduce the comparison    
    next         * PostingNode     // points to the next posting node
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
    }
    defer file.Close()

    // Read file line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // Tokenize the line 
        words := strings.Fields(scanner.Text())
        for _, word := range words {
            word = strings.ToLower(strings.Trim(word, ".,!?\"'"))
            if _ , exists := localIndex[word]; !exists {
                headOfPositions := extractPositions(word, file)
                localIndex[word] = &PostingNode{
                    docId: entry.Name(),
                    positions: headOfPositions,
                    next: nil,
                    jump: nil,
                }
            }
        }

        // Check if there was a problem reading the file
        if err := scanner.Err(); err!= nil {
            log.Print("Error reading file", entry.Name())
        }

    
    }
    // Send the local index to the main goroutine
    indexChan <- localIndex
}


// extract the position for each word in the document
func extractPositions(word string, file *os.File) *PositionNode {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) 

	var head *PositionNode
	var tail *PositionNode

	wordCount := 0 
	for scanner.Scan() {
		wordCount++ 
		currentWord := scanner.Text()

		// check if the current word matches the target word
		if currentWord == word {
			// Add position to linked list
			newNode := &PositionNode{
                position: int16(wordCount),
                next: nil,
                jump: nil,
            }
			if head == nil {
				head = newNode
				tail = newNode
			} else {
				tail.next = newNode
				tail = newNode
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %v\n", err)
		return nil
	}

	return head
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
				for current.next != nil {
					current = current.next
				}
				current.next = postingList
			}
		}
	}

	return index, nil

}

