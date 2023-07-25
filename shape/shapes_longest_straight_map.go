package shape

type ShapesLongestStraightMap struct {
	s       [MAX_NUMBER_OF_CUBES][MAX_NUMBER_OF_CUBES]map[Score]*Shape // array of size, array of longest straight, map of score to pointer to shape
	maxSize ShapeSize
}

func NewShapesLongestStraightMap() Shapes {
	result := &ShapesLongestStraightMap{}
	for i := ShapeSize(0); i < MAX_NUMBER_OF_CUBES; i++ {
		for j := ShapeSize(0); j < MAX_NUMBER_OF_CUBES; j++ {
			result.s[i][j] = make(map[Score]*Shape) // how many shapes do we expect?
		}
	}

	return result
}

func (s ShapesLongestStraightMap) Len() int {
	result := 0
	for i := ShapeSize(0); i < MAX_NUMBER_OF_CUBES; i++ {
		for j := ShapeSize(0); j < MAX_NUMBER_OF_CUBES; j++ {
			result += len(s.s[i][j])
		}
	}

	return result
}

func (s *ShapesLongestStraightMap) Add(shape Shape) Shapes {
	shapeSize := shape.Size()
	longesStraight := shape.LongestStraight()

	s.s[shapeSize-1][longesStraight-1][shape.Score()] = &shape
	if shapeSize > s.maxSize {
		s.maxSize = shapeSize
	}
	return s
}

func (s ShapesLongestStraightMap) GetAll() map[Score]*Shape {
	result := make(map[Score]*Shape)

	for sizeMinOne := range s.s {
		for score, shape := range s.GetAllWithSize(ShapeSize(sizeMinOne + 1)) {
			result[score] = shape
		}
	}

	return result
}

func (s ShapesLongestStraightMap) GetAllWithSize(size ShapeSize) map[Score]*Shape {
	result := make(map[Score]*Shape)

	if size < 1 || size > MAX_NUMBER_OF_CUBES {
		return result
	}

	for _, m := range s.s[size-1] {
		for score, shape := range m {
			result[score] = shape
		}
	}

	return result
}

func (s *ShapesLongestStraightMap) Merge(other Shapes) Shapes {
	for _, shape := range other.GetAll() {
		s.Add(*shape)
	}

	return s
}

func (s ShapesLongestStraightMap) MaxSize() ShapeSize {
	return s.maxSize
}
