// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/lapack"
)

type Dlascler interface {
	Dlascl(kind lapack.MatrixType, kl, ku int, cfrom, cto float64, m, n int, a []float64, lda int)
}

func DlasclTest(t *testing.T, impl Dlascler) {
	const tol = 1e-15

	rnd := rand.New(rand.NewPCG(1, 1))
	for ti, test := range []struct {
		m, n int
	}{
		{0, 0},
		{1, 1},
		{1, 10},
		{10, 1},
		{2, 2},
		{2, 11},
		{11, 2},
		{3, 3},
		{3, 11},
		{11, 3},
		{11, 11},
		{11, 100},
		{100, 11},
	} {
		m := test.m
		n := test.n
		for _, extra := range []int{0, 11} {
			for _, kind := range []lapack.MatrixType{lapack.General, lapack.UpperTri, lapack.LowerTri} {
				a := randomGeneral(m, n, n+extra, rnd)
				aCopy := cloneGeneral(a)
				cfrom := rnd.NormFloat64()
				cto := rnd.NormFloat64()
				scale := cto / cfrom

				impl.Dlascl(kind, -1, -1, cfrom, cto, m, n, a.Data, a.Stride)

				prefix := fmt.Sprintf("Case #%v: kind=%v,m=%v,n=%v,extra=%v", ti, kind, m, n, extra)
				if !generalOutsideAllNaN(a) {
					t.Errorf("%v: out-of-range write to A", prefix)
				}
				switch kind {
				case lapack.UpperTri:
					var mod bool
				loopLower:
					for i := 0; i < m; i++ {
						for j := 0; j < min(i, n); j++ {
							if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
								mod = true
								break loopLower
							}
						}
					}
					if mod {
						t.Errorf("%v: unexpected modification in lower triangle of A", prefix)
					}
				case lapack.LowerTri:
					var mod bool
				loopUpper:
					for i := 0; i < m; i++ {
						for j := i + 1; j < n; j++ {
							if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
								mod = true
								break loopUpper
							}
						}
					}
					if mod {
						t.Errorf("%v: unexpected modification in upper triangle of A", prefix)
					}
				}

				var resid float64
				switch kind {
				case lapack.General:
					for i := 0; i < m; i++ {
						for j := 0; j < n; j++ {
							want := scale * aCopy.Data[i*aCopy.Stride+j]
							got := a.Data[i*a.Stride+j]
							resid = math.Max(resid, math.Abs(want-got))
						}
					}
				case lapack.UpperTri:
					for i := 0; i < m; i++ {
						for j := i; j < n; j++ {
							want := scale * aCopy.Data[i*aCopy.Stride+j]
							got := a.Data[i*a.Stride+j]
							resid = math.Max(resid, math.Abs(want-got))
						}
					}
				case lapack.LowerTri:
					for i := 0; i < m; i++ {
						for j := 0; j <= min(i, n-1); j++ {
							want := scale * aCopy.Data[i*aCopy.Stride+j]
							got := a.Data[i*a.Stride+j]
							resid = math.Max(resid, math.Abs(want-got))
						}
					}
				}
				if resid > tol*float64(max(m, n)) {
					t.Errorf("%v: unexpected result; residual=%v, want<=%v", prefix, resid, tol*float64(max(m, n)))
				}
			}
		}
	}
}
