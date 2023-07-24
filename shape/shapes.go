package shape

const (
	HASH_SIZE = 17
)

type (
	Hash   [HASH_SIZE]uint64
	Shapes struct {
		s       map[int]map[Hash]*Shape
		maxSize int
	}
)

func NewShapes() *Shapes {
	return &Shapes{s: make(map[int]map[Hash]*Shape)}
}

func (s Shapes) Len() int {
	result := 0
	for size := range s.s {
		result += s.NumberOfShapesWithSize(size)
	}

	return result
}

func (s *Shapes) NumberOfShapesWithSize(length int) int {
	return len(s.s[length])
}

func (s *Shapes) Add(shape *Shape) *Shapes {
	shapeSize := shape.Size()
	if _, ok := s.s[shapeSize]; !ok {
		s.s[shapeSize] = make(map[Hash]*Shape)
		if shapeSize > s.maxSize {
			s.maxSize = shapeSize
		}
	}
	s.s[shapeSize][shape.Hash()] = shape

	return s
}

func (s *Shapes) Exists(shape *Shape) bool {
	shapeSize := shape.Size()
	if _, ok := s.s[shapeSize]; !ok {
		return false
	}

	_, ok := s.s[shapeSize][shape.Hash()]
	return ok
}

func (s *Shapes) GetAllWithSize(size int) map[Hash]*Shape {
	if _, ok := s.s[size]; !ok {
		return map[Hash]*Shape{}
	}

	return s.s[size]
}

func (s *Shapes) Merge(other *Shapes) *Shapes {
	for size := range other.s {
		for shapeString := range other.s[size] {
			s.Add(other.s[size][shapeString])
		}
	}
	return s
}

func (s *Shapes) MaxSize() int {
	return s.maxSize
}
