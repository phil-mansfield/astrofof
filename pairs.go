package symfof

import (
	"slices"
)

// Pairer is a buffer struct that holds internal buffers which allow for calls
// to pair-counting functions without excess allocations. These buffers are
// used as return values, meaning that the same Pairer can't have 
type Pairer struct {
	// If set to true, Pairer stops as soon as it has found a single pair.
	StopEarly bool
	// All of these are internal buffers that are meaningless to users.
	i1, i2 []int64
}

func (pair *Pairer) FindPairsOneCell(
	p []Particle, r float32, sortDim int64,
) (i1, i2 []int64) {
	r2 := r*r
	pair.i1, pair.i2 = pair.i1[:0], pair.i2[:0]
	if sortDim == -1 {
		for i := 0; i < len(p) - 1; i++ {
			for j := i + 1; j < len(p); j++ {
				dx := p[i].X[0] - p[j].X[0]
				dy := p[i].X[1] - p[j].X[1]
				dz := p[i].X[2] - p[j].X[2]

				dr2 := dx*dx + dy*dy + dz*dz

				if dr2 <= r2 {
					pair.i1 = append(pair.i1, int64(i))
					pair.i2 = append(pair.i2, int64(j))
				}
			}
		}
	} else {
		high := 1
		for i := 0; i < len(p)-1; i++ {
			for ; high < len(p); high++ {
				delta := p[high].X[sortDim] - p[i].X[sortDim]
				if delta > r { break }
			}

			for j := i+1; j < high; j++ {
				dx := p[i].X[0] - p[j].X[0]
				dy := p[i].X[1] - p[j].X[1]
				dz := p[i].X[2] - p[j].X[2]

				dr2 := dx*dx + dy*dy + dz*dz

				if dr2 <= r2 {
					pair.i1 = append(pair.i1, int64(i))
					pair.i2 = append(pair.i2, int64(j))
				}
			}
		}
	}
	return pair.i1, pair.i2
}

func (pair *Pairer) FindPairsTwoCells(
	p1, p2 []Particle, r float32, sortDim int64,
) (i1, i2 []int64) {
	r2 := r*r
	pair.i1, pair.i2 = pair.i1[:0], pair.i2[:0]
	if sortDim == -1 {
		for i := range p1 {
			for j := range p2 {
				dx := p1[i].X[0] - p2[j].X[0]
				dy := p1[i].X[1] - p2[j].X[1]
				dz := p1[i].X[2] - p2[j].X[2]

				dr2 := dx*dx + dy*dy + dz*dz

				if dr2 <= r2 {
					pair.i1 = append(pair.i1, int64(i))
					pair.i2 = append(pair.i2, int64(j))
				}
			}
		}
	} else {
		low, high := 0, 0
		for i := range p1 {
			for ; low < len(p2); low++ {
				delta := p1[i].X[sortDim] - p2[low].X[sortDim]
				if delta <= r { break }
			}
			for ; high < len(p2); high++ {
				delta := p2[high].X[sortDim] - p1[i].X[sortDim]
				if delta > r { break }
			}

			for j := low; j < high; j++ {
				dx := p1[i].X[0] - p2[j].X[0]
				dy := p1[i].X[1] - p2[j].X[1]
				dz := p1[i].X[2] - p2[j].X[2]

				dr2 := dx*dx + dy*dy + dz*dz

				if dr2 <= r2 {
					pair.i1 = append(pair.i1, int64(i))
					pair.i2 = append(pair.i2, int64(j))
				}
			}
		}
	}
	return pair.i1, pair.i2
}


func (pair *Pairer) SortParticles(p []Particle, dim int) {
	switch len(p) {
	case 0, 1:
	case 2:
		if p[0].X[dim] > p[1].X[dim] {
			p[0], p[1] = p[1], p[0]
		}
	case 3:
		maxi, midi, mini := sort3Index(
			p[0].X[dim], p[1].X[dim], p[2].X[dim], 0, 1, 2,
		)
		p[0], p[1], p[2] = p[mini], p[midi], p[maxi]
	default:
		switch dim {
		case 0: slices.SortFunc(p, ParticleXCmp)
		case 1:	slices.SortFunc(p, ParticleYCmp)
		case 2:	slices.SortFunc(p, ParticleZCmp)
		}
	}
}

func sort3Index(x, y, z float32, ix, iy, iz int) (maxi, midi, mini int) {
	if x > y {
		if x > z {
			if y > z {
				return ix, iy, iz
			} else {
				return ix, iz, iy
			}
		} else {
			return iz, ix, iy
		}
	} else {
		if y > z {
			if x > z {
				return iy, ix, iz
			} else {
				return iy, iz, ix
			}
		} else {
			return iz, iy, ix
		}
	}
}