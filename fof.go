package symfof

func FOF(L float32, x, cen [][3]float32, r float32, nGrid, nMin int) (groups, cenGroups []int32) {
	f := NewFinder(L, x, nGrid)
	uf := NewUnionFinder(int32(len(x)))

	for i := int32(0); i < int32(len(x)); i++ {
		idx := f.Find(x[i], r)
		for _, j := range idx {
			if i == j { continue }
			uf.Union(i, j)
		}
	}

	groups = make([]int32, len(x))
	for i := range x {
		groups[i] = uf.Find(int32(i))
		if uf.Size[groups[i]] < int32(nMin) {
			groups[i] = -1
		}
	}

	cenGroups = make([]int32, len(cen))
	for i := range cenGroups {
		idx := f.Find(cen[i], r)
		if len(idx) == 0 {
			cenGroups[i] = -1
		} else {
			cenGroups[i] = groups[idx[0]]
		}
	}

	return groups, cenGroups
}
