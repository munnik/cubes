package shape

// Shapes is a map of Shape.Len() to a map of Score.Hash() to *Shape
type Shapes map[int]map[string]*Shape

func NewShapes() Shapes {
	return make(Shapes)
}

func (s Shapes) Len() int {
	result := 0
	for index := range s {
		result += len(s[index])
	}

	return result
}

func (s *Shapes) Add(shape *Shape) *Shapes {
	if _, ok := (*s)[shape.Len()]; !ok {
		(*s)[shape.Len()] = make(map[string]*Shape)
	}
	(*s)[shape.Len()][shape.String()] = shape

	return s
}

func (s *Shapes) Merge(other Shapes) *Shapes {
	for len := range other {
		for hash := range other[len] {
			s.Add(other[len][hash])
		}
	}
	return s
}
