package utils

import (
	"errors"
	"math"

	"gonum.org/v1/gonum/mat"
)


func CosineSimilarity(vec1, vec2 []float64) (float64, error) {
	// check if both vectors have the same length
	if len(vec1) != len(vec2) {
		return 0, errors.New("vectors must have the same length")
	}

	// initialize variables for the dot product and magnitudes
	dotProduct := 0.0
	magnitudeVec1 := 0.0
	magnitudeVec2 := 0.0

	// iterate over the vectors to calculate the components
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		magnitudeVec1 += vec1[i] * vec1[i]
		magnitudeVec2 += vec2[i] * vec2[i]
	}

	// compute magnitudes
	magnitudeVec1 = math.Sqrt(magnitudeVec1)
	magnitudeVec2 = math.Sqrt(magnitudeVec2)

	// avoid division by zero
	if magnitudeVec1 == 0 || magnitudeVec2 == 0 {
		return 0, errors.New("one or both vectors have zero magnitude")
	}

	// calculate cosine similarity
	cosineSimilarity := dotProduct / (magnitudeVec1 * magnitudeVec2)
	return cosineSimilarity, nil
}


func LSIcosineSimilarity(a, b *mat.Dense) float64 {
	// Extract the raw data slices from matrices
	aVec := a.RawMatrix().Data
	bVec := b.RawMatrix().Data

	// Dot product
	dotProduct := 0.0
	for i := 0; i < len(aVec); i++ {
		dotProduct += aVec[i] * bVec[i]
	}

	// Norms (lengths) of the vectors
	normA := 0.0
	for i := 0; i < len(aVec); i++ {
		normA += aVec[i] * aVec[i]
	}
	normA = math.Sqrt(normA)

	normB := 0.0
	for i := 0; i < len(bVec); i++ {
		normB += bVec[i] * bVec[i]
	}
	normB = math.Sqrt(normB)

	// Cosine similarity
	return dotProduct / (normA * normB)
}


func CalculateSimilarities(query * mat.Dense , Vt * mat.Dense, k int) []float64 {
	// Step 2: Get number of documents
	_, numDocs := Vt.Dims()

	// Step 3: Calculate cosine similarity with each document in V^T
	similarities := make([]float64, numDocs)
	for i := 0; i < numDocs; i++ {
		// Extract the i-th document from V^T (it is a column in V^T)
		docVec := mat.NewDense(k, 1, nil)
		for j := 0; j < k; j++ {
			docVec.Set(j, 0, Vt.At(j, i))
		}

		// Calculate cosine similarity between query and this document
		similarities[i] = LSIcosineSimilarity(query, docVec)
	}

	return similarities
}
