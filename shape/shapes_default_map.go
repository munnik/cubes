package shape

type ShapesDefaultMap struct {
	s       map[ShapeSize]map[Score]*Shape
	maxSize ShapeSize
}

func NewShapesDefaultMap() Shapes {
	return &ShapesDefaultMap{s: make(map[ShapeSize]map[Score]*Shape)}
}

func (s ShapesDefaultMap) Len() int {
	result := 0
	for size := range s.s {
		result += len(s.s[size])
	}

	return result
}

func (s *ShapesDefaultMap) Add(shape Shape) Shapes {
	shapeSize := shape.Size()
	if _, ok := s.s[shapeSize]; !ok {
		s.s[shapeSize] = make(map[Score]*Shape)
		if shapeSize > s.maxSize {
			s.maxSize = shapeSize
		}
	}
	s.s[shapeSize][shape.Score()] = &shape

	return s
}

func (s ShapesDefaultMap) GetAll() map[Score]*Shape {
	result := make(map[Score]*Shape)

	for size := range s.s {
		for score, shape := range s.GetAllWithSize(size) {
			result[score] = shape
		}
	}

	return result
}

func (s ShapesDefaultMap) GetAllWithSize(size ShapeSize) map[Score]*Shape {
	if _, ok := s.s[size]; !ok {
		return map[Score]*Shape{}
	}

	return s.s[size]
}

func (s *ShapesDefaultMap) Merge(other Shapes) Shapes {
	for _, shape := range other.GetAll() {
		s.Add(*shape)
	}

	return s
}

func (s ShapesDefaultMap) MaxSize() ShapeSize {
	return s.maxSize
}
