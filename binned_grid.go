package symfof

import (
	"slices"
)

// BinnedGrid is a 3d grid which can return an array of items in each grid
// cell. This is designed to be high performance, which means the following:
//
// (1) The input particles are in code units.
// 
// (2) It assumes the origin of the points is zero.
//
// (3) It does not handle periodic checks. If the particles are from a 
// periodic box, they need to have been shifted into their most contiguous
// frame and must be smaller than grid minus a single grid cell.
//
// (4) The upper bounf of the grid is exclusive, not inclusive. The lower bound
// is inclusive.
//
// All particles need to be inserted at once.
type BinnedGrid interface {
	// MakeGrid initializes a grid  with (x, y, z) widths given by span. If
	// called repeatedly, it will dynamically update internal state as much as
	// it can to accomodate any changes in the total number of grid cells being
	// requested.
	Resize(span [3]int64)
	// Bin bins the particles in the array p. Calling a second time clears the
	// bins, since most of these methods can't do dynamic insertions. Particles
	// may be rearranged within p.
	Bin(p []Particle)
	// IsEmpty returns true if there are no particles at the specified bin
	// and false otherwise.
	// Size returns the number of elements in the specified bin.
	Size(idx [3]int64) int64
	// Get returns all the particles in a bin in no particular order. For some
	// binning schemes, this will always be an allocation-free operation. For
	// other schemes, particles are stored in non-array data structures. In
	// these cases an optional out array can be supplied as the first variadic
	// argument. out may be resized by Get, so you should re-assign it as the
	// return value of Get instead of continuing to reference the same array.
	// []Particle arrays beyond the first are ignored.
	Get(idx [3]int64, out ...[]Particle) []Particle
}

/////////////////////////////
// ArrayListGrid functions //
/////////////////////////////

// ArrayListGrid implements the BinnedGrid interface using individual
// dynamically allocated arrays (i.e. Go slices) in each bin.
type ArrayListGrid struct {
	// Data is a flat array representing the grid. It uses C ordering, so
	// the x is the fastest changing dimension, then y, then z.
	Data [][]Particle
	// Span is the dimensions of the grid.
	Span [3]int64
	// Dy, and Dz are convenience parameters to make grid math easier and
	// represent how many cells along the flat array need to be traveled to
	// increment a given coordinate by one.
	Dy, Dz int64
}

func (g *ArrayListGrid) Resize(span [3]int64) {
	g.Span = span
	g.Dy, g.Dz = span[0], span[0]*span[1]
	n := span[0]*span[1]*span[2]

	g.Data = slices.Grow(g.Data, int(n))[:n]

	// After much internal debate, I've decided not to just slice the arrays
	// down to zero. I think that will leads to ballooning memory costs.
	for i := range g.Data { g.Data[i] = nil }
}

func (g *ArrayListGrid) Bin(p []Particle) {
	for i := range p {
		pi := p[i]
		ix, iy, iz := int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])
		j := ix + iy*g.Dy + iz*g.Dz
		g.Data[j] = append(g.Data[j], pi)
	}
}

func (g *ArrayListGrid) Size(idx [3]int64) int64 {
	j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz	
	return int64(len(g.Data[j]))
}

// The Get method of ArrayListGrid can return an array without making new
// allocations and does not need a buffer.
func (g *ArrayListGrid) Get (idx [3]int64, out ...[]Particle) []Particle {
	j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz	
	return g.Data[j]
}
var _ BinnedGrid = &ArrayListGrid{ }

///////////////////////////////////
// NaiveLinkedListGrid functions //
///////////////////////////////////

// NaiveLinkedListGrid implements the BinnedGrid interface using linked lists
// where each element is separately allocated.
type NaiveLinkedListGrid struct {
	// Heads is a flat array representing a grid with C-ordered indices that
	// stores and stores the first node in each list. May be nil.
	Heads []*ListNode
	// Sizes is the length the lists.
	Sizes []int64
	// Dy, and Dz are convenience parameters to make grid math easier and
	// represent how many cells along the flat array need to be traveled to
	// increment a given coordinate by one.
	Dy, Dz int64
}

// ListNode is a single node in a pointer-based linked list.
type ListNode struct {
	P Particle
	Next *ListNode
}

func (g *NaiveLinkedListGrid) Resize(span [3]int64) {
	n := span[0]*span[1]*span[2]
	g.Dy, g.Dz = span[0], span[0]*span[1]

	g.Heads = slices.Grow(g.Heads, int(n))[:n]
	g.Sizes = slices.Grow(g.Sizes, int(n))[:n]
	for i := range g.Sizes {
		g.Heads[i], g.Sizes[i] = nil, 0
	}
}

func (g *NaiveLinkedListGrid) Bin(p []Particle) {
	for i := range p {
		pi := p[i]
		idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
		j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
		node := &ListNode{ pi, g.Heads[j] }
		g.Heads[j] = node
		g.Sizes[j]++
	}
}

func (g *NaiveLinkedListGrid) Size(idx [3]int64) int64 {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	return g.Sizes[i]
}
// The Get method of NaiveLinkedListGrid will need to make internal allocations
// unless passed a buffer.
func (g *NaiveLinkedListGrid) Get (idx [3]int64, out ...[]Particle) []Particle {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	n := g.Sizes[i]

	var buf []Particle
	if len(out) == 0 {
		buf = make([]Particle, n)
	} else {
		buf = slices.Grow(out[0], int(n))[:n]
	}

	node := g.Heads[i]
	for i := range buf {
		buf[i] = node.P
		node = node.Next
	}

	return buf
}
var _ BinnedGrid = &NaiveLinkedListGrid{ }

//////////////////////////////
// LinkedListGrid functions //
//////////////////////////////

// LinkedListGrid implements the BinnedGrid interface using a flat array of
// indices.
type LinkedListGrid struct {
	// Heads is a flat array representing a grid with C-ordered indices that
	// stores and stores the first node in each list. May be nil.
	Heads []int64
	// Sizes is the length the lists.
	Sizes []int64
	// Dy, and Dz are convenience parameters to make grid math easier and
	// represent how many cells along the flat array need to be traveled to
	// increment a given coordinate by one.
	Dy, Dz int64

	// The particles are represented by flat arrays, not indiviudal nodes
	// as is used in the naive case.
	Data []Particle
	Next []int64
}

func (g *LinkedListGrid) Resize(span [3]int64) {
	n := span[0]*span[1]*span[2]
	g.Dy, g.Dz = span[0], span[0]*span[1]

	g.Heads = slices.Grow(g.Heads, int(n))[:n]
	g.Sizes = slices.Grow(g.Sizes, int(n))[:n]
	for i := range g.Sizes {
		g.Heads[i], g.Sizes[i] = -1, 0
	}

	g.Data = nil
	g.Next = g.Next[:0]
}

func (g *LinkedListGrid) Bin(p []Particle) {
	g.Next = slices.Grow(g.Next, len(p))[:len(p)]
	g.Data = p

	for i := range p {
		pi := p[i]
		idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
		j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
		
		g.Next[i] = g.Heads[j]
		g.Heads[j] = int64(i)
		g.Sizes[j]++
	}
}

func (g *LinkedListGrid) Size (idx [3]int64) int64 {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	return g.Sizes[i]
}

// The Get method of LinkedListGrid will need to make internal allocations
// unless passed a buffer.
func (g *LinkedListGrid) Get (idx [3]int64, out ...[]Particle) []Particle {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	n := g.Sizes[i]

	var buf []Particle
	if len(out) == 0 {
		buf = make([]Particle, n)
	} else {
		buf = slices.Grow(out[0], int(n))[:n]
	}

	node := g.Heads[i]
	for i := range buf {
		buf[i] = g.Data[node]
		node = g.Next[node]
	}

	return buf
}

var _ BinnedGrid = &LinkedListGrid{ }

////////////////////////////////
// CountingSortGrid functions //
////////////////////////////////

// CountingSortGrid implements the BinnedGrid interface through the counting
// sort algorithm.
type CountingSortGrid struct {
	Counts, BinEdges []int64 // Keep track of bin sizes, slices of same array
	Data []Particle // Sorted data is an internal buffer

	// Dy, and Dz are convenience parameters to make grid math easier and
	// represent how many cells along the flat array need to be traveled to
	// increment a given coordinate by one.
	Dy, Dz int64
}

func (g *CountingSortGrid) Resize(span [3]int64) {
	g.Dy, g.Dz = span[0], span[1]*span[0]
	n := span[0]*span[1]*span[2]

	g.BinEdges = slices.Grow(g.BinEdges, int(n)+1)[:n+1]
	g.Counts = g.BinEdges[1:]
	for i := range g.BinEdges { g.BinEdges[i] = 0 }
}

func (g *CountingSortGrid) Bin(p []Particle) {
	g.Data = slices.Grow(g.Data, len(p))[:len(p)]

	// Count particles first
	for _, pi := range p {
		idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
		j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
		g.Counts[j]++
	}

	// Turn the g.Counts array into an array of bin starts. You need to swap
	// things between currCount and prevCount becasue these are all slices
	// of the same array.
	binStarts := g.Counts
	prevCount := int64(0)
	for i := range binStarts {
		currCount := g.Counts[i]
		binStarts[i] = g.BinEdges[i] + prevCount
		prevCount = currCount
	}

	// Add particles into their bins, keeping track of everything with
	// binStarts. This ensures g.BinEdges has the correct values.
	for _, pi := range p {
		idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
		j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
		g.Data[binStarts[j]] = pi
		binStarts[j]++
	}
}

func (g *CountingSortGrid) Size(idx [3]int64) int64 {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	return g.BinEdges[i+1] - g.BinEdges[i]
}

// The Get method of ArrayListGrid can return an array without making new
// allocations and does not need a buffer.
func (g *CountingSortGrid) Get (idx [3]int64, out ...[]Particle) []Particle {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	return g.Data[g.BinEdges[i]: g.BinEdges[i+1]]
}

var _ BinnedGrid = &CountingSortGrid{ }

//////////////////////////////
// CycleSorttGrid functions //
//////////////////////////////

// CycleSortGrid implements the BinnedGrid interface through the cycle sort
// algorithm.
type CycleSortGrid struct {
	Counts, BinEdges, BinEnds []int64
	Data []Particle // Sorted data IS NOT an internal buffer!

	// Dy, and Dz are convenience parameters to make grid math easier and
	// represent how many cells along the flat array need to be traveled to
	// increment a given coordinate by one.
	Dy, Dz int64
}
func (g *CycleSortGrid) Resize(span [3]int64) {
	g.Dy, g.Dz = span[0], span[1]*span[0]
	n := span[0]*span[1]*span[2]

	g.BinEdges = slices.Grow(g.BinEdges, int(n)+1)[:n+1]
	g.Counts = g.BinEdges[1:]
	g.BinEnds = slices.Grow(g.BinEnds, int(n))[:n]
	for i := range g.BinEdges { g.BinEdges[i] = 0 }
}
func (g *CycleSortGrid) Bin(p []Particle) {
	g.Data = p
	// Count particles first
	for _, pi := range p {
		idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
		j := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
		g.Counts[j]++
	}

	for i := 1; i < len(g.BinEdges); i++ {
		g.BinEdges[i] += g.BinEdges[i-1]
	}
	for i := range g.BinEnds {
		g.BinEnds[i] = g.BinEdges[i]
	}

	// Copied this from an old blog post I wrote years ago. Yes, I also think
	// this is super complicated. But keep reading it and you'll get it.
    for srcBin := int64(0); srcBin < int64(len(g.BinEnds)); srcBin++ {
        for lim := g.BinEdges[srcBin+1]; g.BinEnds[srcBin] < lim; {
            // i is the index of element we're going to correct the position of.
            i := g.BinEnds[srcBin]

            pi := p[i]
            idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
			dstBin := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
            for ; dstBin != srcBin; {
                // j is the index of the element we're kicking out.
                j := g.BinEnds[dstBin]
                //fmt.Printf("   len(p) %d, src %d, dst %d, i %d, j %d\n",
                //	len(p), srcBin, dstBin, i, j)
                //fmt.Printf("%d\n", g.BinEnds)
                p[i], p[j] = p[j], p[i]
                g.BinEnds[dstBin]++

                pi = p[i]
				idx := [3]int64{int64(pi.X[0]), int64(pi.X[1]), int64(pi.X[2])}
				dstBin = idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
            }

            g.BinEnds[srcBin]++
        }
    }
}

func (g *CycleSortGrid) Size(idx [3]int64) int64 {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	return g.BinEdges[i+1] - g.BinEdges[i]
}
// The Get method of ArrayListGrid can return an array without making new
// allocations and does not need a buffer.
func (g *CycleSortGrid) Get (idx [3]int64, out ...[]Particle) []Particle {
	i := idx[0] + idx[1]*g.Dy + idx[2]*g.Dz
	return g.Data[g.BinEdges[i]: g.BinEdges[i+1]]
}
var _ BinnedGrid = &CycleSortGrid{ }

