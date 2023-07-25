package shape

type ShapesMap struct {
	s       map[ShapeSize]map[Score]*Shape
	maxSize ShapeSize
}

func NewShapesMap() *ShapesMap {
	return &ShapesMap{s: make(map[ShapeSize]map[Score]*Shape)}
}

func (s ShapesMap) Len() int {
	result := 0
	for size := range s.s {
		result += s.NumberOfShapesMapWithSize(size)
	}

	return result
}

func (s *ShapesMap) NumberOfShapesMapWithSize(size ShapeSize) int {
	return len(s.s[size])
}

func (s *ShapesMap) Add(shape Shape) Shapes {
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

func (s *ShapesMap) GetAll() map[Score]*Shape {
	result := make(map[Score]*Shape)

	for size := range s.s {
		for score, shape := range s.GetAllWithSize(size) {
			result[score] = shape
		}
	}

	return result
}

func (s *ShapesMap) GetAllWithSize(size ShapeSize) map[Score]*Shape {
	if _, ok := s.s[size]; !ok {
		return map[Score]*Shape{}
	}

	return s.s[size]
}

func (s *ShapesMap) Merge(other Shapes) Shapes {
	for _, shape := range other.GetAll() {
		s.Add(*shape)
	}

	return s
}

func (s *ShapesMap) MaxSize() ShapeSize {
	return s.maxSize
}
