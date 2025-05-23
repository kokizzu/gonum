// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
)

func DtrmvBenchmark(b *testing.B, dtrmv Dtrmver, n, lda, incX int, ul blas.Uplo, tA blas.Transpose, d blas.Diag) {
	rnd := rand.New(rand.NewPCG(0, 0))
	a := make([]float64, n*lda)
	for i := range a {
		a[i] = rnd.Float64()
	}

	x := make([]float64, n*incX)
	for i := range x {
		x[i] = rnd.Float64()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dtrmv.Dtrmv(ul, tA, d, n, a, lda, x, incX)
	}
}
