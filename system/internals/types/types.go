package types

// TDM with TF-IDF weights
type TDM struct {
	Matrix      map[string]map[string]float64 // Term -> Document -> TF-IDF Weight
	Terms       []string                      // List of all terms
	Documents   []string                      // List of all documents
	DocWordCount map[string]int               // Total word count for each document
}