package symfof

import (
	"testing"
)

//////////////////////
// CorrectnessTests //
//////////////////////

var gridModels = []BinnedGrid{
	new(ArrayListGrid),
	new(NaiveLinkedListGrid),
	new(LinkedListGrid),
	new(CountingSortGrid),
	new(CycleSortGrid),
}
var gridModelNames = []string{
	"ArrayListGrid",
	"NaiveLinkedListGrid",
	"LinkedListGrid",
	"CountngSortGrid",
	"CycleSortGrid",
}

// TestSmallGrid inserts a small number of points into a small grid and checks
// that the binning returns exactly the grid you expect.
func TestSmallGrid(t *testing.T) {
	span := [3]int64{ 2, 3, 2 }
	zero := [3]float32{}
	points := []Particle {
		{0, [3]float32{0.5, 0.5, 0.5}, zero},
		{1, [3]float32{0.55, 0.25, 0.15}, zero},
		{2, [3]float32{1.5, 2.0, 0.0}, zero},
		{30000, [3]float32{1.2, 1.9, 1.1}, zero},
		{3, [3]float32{1, 2, 0}, zero},
		{10, [3]float32{1, 1, 1}, zero},
		{9, [3]float32{1.999, 2.999, 0.999}, zero},
	}
	cells := [][3]int64{
		{0, 0, 0}, {0, 0, 0}, {1, 2, 0}, {1, 1, 1},
		{1, 2, 0}, {1, 1, 1}, {1, 2, 0},
	}
	ids := make([]uint64, len(points))
	for i := range points { ids[i] = points[i].ID }

	for ig := range gridModels {
		grid, gridName := gridModels[ig], gridModelNames[ig]
		
		grid.Resize(span)
		grid.Bin(points)

		buf := []Particle{ }

		totalSize := 0

		for iz := int64(0); iz < span[2]; iz++ {
			for iy := int64(0); iy < span[1]; iy++ {
				for ix := int64(0); ix < span[0]; ix++ {
					idx := [3]int64{ix, iy, iz}
					p := grid.Get(idx, buf)
					if grid.Size(idx) != int64(len(p)) {
						t.Errorf("%s: Size() of cell %d is %d but should " + 
							"be %d.", gridName, idx, grid.Size(idx), len(p))
					}

					totalSize += len(p)

					for i := range points {
						if cells[i] == idx {
							foundParticle := false
							for j := range p {
								if p[j].ID == ids[i] {
									foundParticle = true
								}
							}
							if !foundParticle {
								t.Errorf("%s: particle with ID %d was found " + 
									"in cell %d, but shouldn't have been.",
									gridName, points[i].ID, idx)
							}
						}
					}
				}
			}
		}

		// Combined with the previous tests, this ensures that no point is
		// showing up where it shouldn't be.
		if totalSize != len(points) {
			t.Errorf("%s: total number of particles is %d but should " + 
				"have been %d", gridName, totalSize, len(points))
		}
	}
}

// TestResize repeatedly resizes grids up and down, inserts three elements into
// each grid cell and makes sure that three elements are found.
func TestResize(t *testing.T) {
	spans := [][3]int64{
		{0, 0, 0},
		{10, 10, 10},
		{1, 3, 2},
		{0, 0, 0},
		{10, 10, 10},
	}

	points := make([][]Particle, len(spans))

	reps := 3
	for i := range spans {
		p := []Particle{ }
		span := spans[i]
		for r := 0; r < reps; r++ {
			for iz := int64(0); iz < span[2]; iz++ {
				for iy := int64(0); iy < span[1]; iy++ {
					for ix := int64(0); ix < span[0]; ix++ {
						offset := float32(0.5)
						x := [3]float32{
							float32(ix)+offset,
							float32(iy)+offset,
							float32(iz)+offset,
						}
						p = append(p, Particle{0, x, [3]float32{0, 0, 0}})
					}
				}
			}
		}
		points[i] = p
	}

	for ig := range gridModels {
		grid, gridName := gridModels[ig], gridModelNames[ig]
		for is := range spans {
			span := spans[is]
			grid.Resize(span)
			grid.Bin(points[is])
			for iz := int64(0); iz < span[2]; iz++ {
				for iy := int64(0); iy < span[1]; iy++ {
					for ix := int64(0); ix < span[0]; ix++ {
						idx := [3]int64{ix, iy, iz}
						if grid.Size(idx) != int64(reps) {
							t.Errorf("%s, span = %d: cell %d has %d elements " + 
								"but should have %d.", gridName, span, idx,
								grid.Size(idx), reps)
						}
					}
				}
			}
		}
	}
}
