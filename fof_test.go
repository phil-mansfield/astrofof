package symfof

import (
	"testing"
)

func TestFOF(t *testing.T) {
	L := float32(200)
	r := float32(2)
	groupIdxs := [][]int32{
		{0, 1, 2, 3},
		{4, 5, 6},
		{8, 9, 10, 11, 12, 13},
	}
	freeParticles := []int32{ 7, 14, 15 }
	x := [][3]float32{
		{100, 100, 100}, // Group 0
		{101, 100, 100},
		{102, 100, 100},
		{103, 100, 100},

		{150, 199, 20}, // Group 1
		{150, 0, 20},
		{150, 1, 20},

		{1, 2, 3}, // --
		
		{75, 75, 30}, // Group 2
		{75, 75, 29},
		{75, 74, 31},
		{74, 74, 31},
		{73, 74, 31},
		{75, 75, 31},

		{120, 120, 120},
		{120, 120, 120},
	}

	cenGroupIdxs := []int32{ 0, 1, 2, -1 }
	cen := [][3]float32{
		{100.5, 100, 100},
		{150.5, 1.5, 20.5},
		{75, 75, 30},
		{120, 120, 120},
	}

	groups, cenGroups := FOF(L, x, cen, r, 10, 3)

	for _, i := range freeParticles {
		if groups[i] != -1 {
			t.Errorf("Expected particle %d to be in no group, but it is in group %d", i, groups[i])
		}
	}
	
	roots := make([]int32, len(groupIdxs))
	for i := range groupIdxs {
		roots[i] = groups[groupIdxs[i][0]]
	}

	for i, idxs := range groupIdxs {
		for _, j := range idxs {
			if groups[j] != roots[i] {
				t.Errorf("Expected particle %d to have group %d, found group %d", j, roots[i], groups[j])
			}
		}
	}

	for i := range cenGroups {
		if cenGroupIdxs[i] == -1 && cenGroups[i] != -1 {
			t.Errorf("Expected central %d to have grpup %d, fround group %d",
				i, -1, cenGroups[i])
		} else if cenGroupIdxs[i] != -1 && roots[cenGroupIdxs[i]] != roots[i] {
			t.Errorf("Expected central %d to have grpup %d, fround group %d",
				i, roots[cenGroupIdxs[i]], cenGroups[i])
		}
	}
}
