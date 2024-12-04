package utils

import (
	"gonum.org/v1/gonum/mat"
)


func SVD(A [][]float64) (*mat.Dense, *mat.Dense, *mat.Dense) {
	// [][]float64 to *mat.Dense
	rows := len(A)
	cols := len(A[0])

	//log.Println(A)

	flatData := make([]float64, 0, rows*cols)
	for _, row := range A {
		flatData = append(flatData, row...)
	}


	matA := mat.NewDense(rows, cols, flatData)

	//log.Println(matA)

	// Perform SVD decomposition
	var svd mat.SVD
	ok := svd.Factorize(matA, mat.SVDThin) 
	if !ok {
		panic("SVD factorization failed")
	}

	// Extract U, Sigma, and V^T
	var U, VT mat.Dense
	svd.UTo(&U)
	svd.VTo(&VT)

	// Extract singular values and construct Sigma matrix
	sigma := svd.Values(nil)
	sigmaMat := mat.NewDense(len(sigma), len(sigma), nil)
	for i := 0; i < len(sigma); i++ {
		sigmaMat.Set(i, i, sigma[i]) 
	}

	return &U, sigmaMat, &VT
}