package symfof

type UnionFinder struct {
	Parent []int32
	Size []int32
	NGroup int32
}

func NewUnionFinder(n int32) *UnionFinder {
	uf := &UnionFinder{
		Parent: make([]int32, n),
		Size: make([]int32, n),
	}
	for i := range uf.Parent {
		uf.Parent[i] = int32(i)
		uf.Size[i] = 1
	}

	return uf
}

func (uf *UnionFinder) Find(i int32) int32 {
	j := i
	for uf.Parent[j] != j {
		j = uf.Parent[j]
	}
	root :=  j
	for j = i; uf.Parent[j] != j; {
		next := uf.Parent[j]
		uf.Parent[j] = root
		j = next
	}

	return root
}

func (uf *UnionFinder) Union(i, j int32) {
	rooti, rootj := uf.Find(i), uf.Find(j)
	if rooti == rootj { return }
	sizei, sizej := uf.Size[rooti], uf.Size[rootj]
	if sizei < sizej {
		uf.Parent[rooti] = rootj
		uf.Size[rootj] = sizei + sizej
	} else {
		uf.Parent[rootj] = rooti
		uf.Size[rooti] = sizei + sizej
	}
	uf.NGroup--
}
