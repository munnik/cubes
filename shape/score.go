package shape

import (
	"crypto/sha256"
	"encoding/binary"
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

func (s Score) Hash() string {
	b := make([]byte, 8)
	h := sha256.New()
	indices := s.SortIndices()

	for _, index := range indices {
		binary.LittleEndian.PutUint64(b, index)
		h.Write(b)
	}
	return (string)(h.Sum(nil))
}

func (s Score) SortIndices() []uint64 {
	indices := make([]uint64, 0, len(s))
	for index := range s {
		indices = append(indices, index)
	}
	sort.Slice(indices, func(i, j int) bool { return indices[i] < indices[j] })
	return indices
}
