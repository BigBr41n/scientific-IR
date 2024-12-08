package models


func (data * Data) ClassicBoolean(query []string) ([]string, error) {

    // store the set of matching DocIDs for each term.
    var intersectedDocIDs map[string]struct{}
    isFirstTerm := true // handle the first term differently.

	for _, word := range query {
        posting := (*data.InvertedIndex)[word]

		currentDocIDs := make(map[string]struct{})

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