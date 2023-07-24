package shape

import (
	"sort"
)

type Score map[uint64]bool

// Compare s to other, return -1 if left is smaller than right, 0 if left is equal to right and 1 if left is bigger than right
func (left Score) Cmp(right Score) int {
	leftIndices := left.SortIndices()
	rightIndices := right.SortIndices()
	var leftIndex, rightIndex uint64

	for len(leftIndices) > 0 && len(rightIndices) > 0 {
		leftIndex = leftIndices[len(leftIndices)-1]
		rightIndex = rightIndices[len(leftIndices)-1]

		if leftIndex > rightIndex {
			return 1
		}
		if leftIndex < rightIndex {
			return -1
		}

		leftIndices = leftIndices[:len(leftIndices)-1]
		rightIndices = rightIndices[:len(rightIndices)-1]
	}

	if len(leftIndices) > 0 {
		return 1
	}
	if len(rightIndices) > 0 {
		return -1
	}

	return 0
}

func (s Score) SortIndices() []uint64 {
	indices := make([]uint64, 0, len(s))
	for index, value := range s {
		if value {
			indices = append(indices, index)
		}
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })
	return indices
}
