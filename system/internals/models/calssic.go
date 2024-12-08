package models

// param : query slice (normalized && stemmed)
// return : slice of strings as a result and error
func (data * Data) ClassicBoolean(query []string) ([]string, error) {

    // store the set of matching DocIDs for each term.
    var intersectedDocIDs map[string]struct{}
    isFirstTerm := true // handle the first term differently.

    // range over the query slice 
	for _, word := range query {

        // grab the posting (linked list)
        posting := (*data.InvertedIndex)[word]

        // create a map of string keys and empty struct e.g. : [document_id]{}
		currentDocIDs := make(map[string]struct{})

        // loop over the linked list and populate the currentDocIDs map.
        for posting != nil {
            currentDocIDs[posting.DocID] = struct{}{}
            posting = posting.Next
        }

        // intersection with the existing results.
        if isFirstTerm {
            intersectedDocIDs = currentDocIDs
            isFirstTerm = false
        } else {
            for docID := range intersectedDocIDs {
                if _, exists := currentDocIDs[docID]; !exists {
                    delete(intersectedDocIDs, docID)
                }
            }
        }
    }

    // Convert the resulting intersection set to a slice of DocIDs.
    result := make([]string, 0, len(intersectedDocIDs))
    for docID := range intersectedDocIDs {
        result = append(result, docID)
    }

    return result, nil
}