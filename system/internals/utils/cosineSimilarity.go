package utils

import (
	"errors"
	"math"
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