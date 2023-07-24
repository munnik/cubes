package shape

// Shapes is a map of Shape.Len() to a map of Score.Hash() to *Shape
type Shapes struct {
	s       map[int]map[string]*Shape
	maxSize int
}

func NewShapes() *Shapes {
	return &Shapes{s: make(map[int]map[string]*Shape)}
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
		s.s[shapeSize] = make(map[string]*Shape)
		if shapeSize > s.maxSize {
			s.maxSize = shapeSize
		}
	}
	s.s[shapeSize][shape.String()] = shape

	return s
}

func (s *Shapes) GetAllWithSize(size int) map[string]*Shape {
	if _, ok := s.s[size]; !ok {
		return map[string]*Shape{}
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
