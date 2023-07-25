package shape

import (
	"sort"
)

const (
	MAX_NUMBER_OF_CUBES = ShapeSize(17)
)

type Score [MAX_NUMBER_OF_CUBES]uint64

func NewScore(s *Shape) *Score {
	result := &Score{}
	size := int(s.Size())
	sizeSquared := size * size

	index := 0
	for c := range s.AllPositiveCoords().coords {
		result[index] = uint64(c[XAxis]) + uint64(c[YAxis]*size) + uint64(c[ZAxis]*sizeSquared) + 1
		index++
	}

	sort.Sort(result)

	return result
}

func (s Score) Len() int {
	return len(s)
}

func (s Score) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s *Score) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Compare s to other, return -1 if left is smaller than right, 0 if left is equal to right and 1 if left is bigger than right
func (left Score) Cmp(right Score) int {
	for index := ShapeSize(0); index < MAX_NUMBER_OF_CUBES; index++ {
		if left[index] < right[index] {
			return -1
		}
		if left[index] > right[index] {
			return 1
		}
	}
	return 0
}
