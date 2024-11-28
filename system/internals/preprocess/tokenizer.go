package preprocess

import "os"

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
func tokenize(file * os.File) (InvertedIndex , error){}


// extract the position for each word in the document
func extractPositions(wordsSlice string , file os.File) (*[]int) {}


// open each file in the dir and pass it to Tokenize  function 
// to get the local inverted index for each doc and then merge them in one inverted index 
func (dp * Tokenizer) ProcessFiles() (InvertedIndex , error) {}

