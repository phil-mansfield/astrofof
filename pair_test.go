package symfof

import (
	"testing"
)

func edgesEqual(i1, i2, j1, j2 []int64) bool {
	if len(i1) != len(j1) || len(j1) != len(j2) || len(i1) != len(i2) {
	}

	alreadyUsed := make([]bool, len(j1))
iLoop:
	for i := range i1 {
		for j := range j1 {
			if alreadyUsed[j] { continue }
			if (i1[i] == j1[j] && i2[i] == j2[j]) ||
				(i1[i] == j2[j] && i2[i] == j1[j]) {
					alreadyUsed[j] = true
					continue iLoop
			}
		}
		return false
	}
	return true
}

func  TestFindPairsOneCell(t *testing.T) {
	tests := []struct{
		x [][3]float32
		sortDim int64
		i1, i2 []int64
	} {
		// Empty cell, unsorted
		{
			[][3]float32{},
			-1,
			[]int64{},
			[]int64{},
		},
		// Empty cell, sorted
		{
			[][3]float32{},
			0,
			[]int64{},
			[]int64{},
		},
		// Singleton, unsorted
		{
			[][3]float32{{0, 0, 0}},
			-1,
			[]int64{},
			[]int64{},
		},
		// Singleton, sorted
		{
			[][3]float32{{0, 0, 0}},
			1,
			[]int64{},
			[]int64{},
		},
		// Pair, unsorted
		{
			[][3]float32{{0.5, 0.5, 0.5}, {0, 0, 0}},
			-1,
			[]int64{0},
			[]int64{1},
		},
		// Pair, sorted
		{
			[][3]float32{{0, 0, 0}, {-0.5, -0.5, 0.5}},
			2,
			[]int64{0},
			[]int64{1},
		},
		// A bunch of unsorted points
		{
			[][3]float32{{0, 0, 0}, {0.5, 0, 0}, {-0.5, 0, 0}, {0, 0.99, 0},
						{-1, -1, -1}, {-1.5, -1, -1}, {-1, -1.5, -1}},
			-1,
			[]int64{0, 0, 0, 1, 4, 4, 5},
			[]int64{1, 2, 3, 2, 5, 6, 6},
		},
		// A bunch of sorted points
		{
			[][3]float32{
				{-1.5, -1, -1}, {-1, -1.5, -1}, {-1, -1, -1}, {-0.5, 0, 0},
				{0, 0, 0}, {0, 0.99, 0}, {0.5, 0, 0}},
			0,
			[]int64{0, 0, 1, 3, 3, 4, 4},
			[]int64{1, 2, 2, 4, 6, 5, 6},
		},
	}

	pair := &Pairer{ }
	for i := range tests {
		p := make([]Particle, len(tests[i].x))
		for j := range p {
			p[j].X = tests[i].x[j]
		}
		i1, i2 := pair.FindPairsOneCell(p, 1, tests[i].sortDim)
		if !edgesEqual(i1, i2, tests[i].i1, tests[i].i2) {
			t.Errorf("%d) Expected edges %d %d, but got %d %d", i,
				tests[i].i1, tests[i].i2, i1, i2,
			)
		}
	}
}

func  TestFindPairstwoCells(t *testing.T) {
	tests := []struct{
		x1, x2 [][3]float32
		sortDim int64
		i1, i2 []int64
	} {
		// Various combinations of empty and singleton cells
		{ [][3]float32{}, [][3]float32{}, -1, []int64{}, []int64{}, },
		{ [][3]float32{}, [][3]float32{{0,0,0}}, -1, []int64{}, []int64{}, },
		{ [][3]float32{{0,0,0}}, [][3]float32{}, -1, []int64{}, []int64{}, },
		{ [][3]float32{{1,1,1}}, [][3]float32{{0,0,0}}, -1, []int64{}, []int64{}, },
		{ [][3]float32{{0.5,0.5,0.5}}, [][3]float32{{0,0,0}}, -1, []int64{0}, []int64{0}, },
		{ [][3]float32{}, [][3]float32{}, 0, []int64{}, []int64{}, },
		{ [][3]float32{}, [][3]float32{{0,0,0}}, 1, []int64{}, []int64{}, },
		{ [][3]float32{{0,0,0}}, [][3]float32{}, 2, []int64{}, []int64{}, },
		{ [][3]float32{{0,0,0}}, [][3]float32{{1,1,1}}, 0, []int64{}, []int64{}, },
		{ [][3]float32{{0.5,0.5,0.5}}, [][3]float32{{1,1,1}}, 1, []int64{0}, []int64{0}, },


		// A bunch of unsorted points
		{
			[][3]float32{{0,0,0}, {0.5, 0.5, 0.5}, {0,0,1},},
			[][3]float32{{1,1,1}, {1,0,0}, {1,0.5,0.5}, {1.5, 0.5, 0.5}},
			-1,
			[]int64{0, 1, 1, 1, 1},
			[]int64{1, 0, 1, 2, 3},
		},
		// A bunch of sorted points
		{
			[][3]float32{{0,0,1}, {0,0,0}, {0.5, 0.5, 0.5},},
			[][3]float32{{1,1,1}, {1,0,0}, {1,0.5,0.5}, {1.5, 0.5, 0.5}},
			-1,
			[]int64{1, 2, 2, 2, 2},
			[]int64{1, 0, 1, 2, 3},
		},
	}

	pair := &Pairer{ }
	for i := range tests {
		p1 := make([]Particle, len(tests[i].x1))
		p2 := make([]Particle, len(tests[i].x2))
		for j := range p1 {
			p1[j].X = tests[i].x1[j]
		}
		for j := range p2 {
			p2[j].X = tests[i].x2[j]
		}
		i1, i2 := pair.FindPairsTwoCells(p1, p2, 1, tests[i].sortDim)
		if !edgesEqual(i1, i2, tests[i].i1, tests[i].i2) {
			t.Errorf("%d) Expected edges %d %d, but got %d %d", i,
				tests[i].i1, tests[i].i2, i1, i2,
			)
		}
	}
}