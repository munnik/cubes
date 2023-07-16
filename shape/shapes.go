package shape

// Shapes is a map of Shape.Len() to a map of Score.Hash() to *Shape
type Shapes map[int]map[string]*Shape

func (s Shapes) Len() int {
	result := 0
	for index := range s {
		result += len(s[index])
	}

	return result
}

func (s *Shapes) Add(other *Shape) *Shapes {
	if _, ok := (*s)[other.Len()]; !ok {
		(*s)[other.Len()] = make(map[string]*Shape)
	}
	(*s)[other.Len()][other.Score().Hash()] = other

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
