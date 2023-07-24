package shape

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

const SEPARATOR = ", "

type Shape struct {
	coords map[Coord]struct{}
}

func NewShape() *Shape {
	return &Shape{
		coords: map[Coord]struct{}{{0, 0, 0}: {}},
	}
}

// Size is the number of cubes in the collection.
func (s *Shape) Size() int { return len(s.coords) }

func (s *Shape) AddCube(c *Coord) (*Shape, error) {
	if _, ok := s.coords[*c]; ok {
		return nil, fmt.Errorf("coord %v is already in this shape", c)
	}

	if !s.IsNeighbor(c) {
		return nil, fmt.Errorf("coord %v is not a neighbor of this shape", c)
	}

	result := NewShape()
	result.coords = make(map[Coord]struct{}, s.Size()+1)
	for c := range s.coords {
		result.coords[c] = struct{}{}
	}
	result.coords[*c] = struct{}{}

	return result, nil
}

func (s *Shape) MustAddCube(c *Coord) *Shape {
	result, err := s.AddCube(c)
	if err != nil {
		panic(err)
	}

	return result
}

// returns all possible new shapes with one cube added to the original shape
func (s *Shape) Grow() *Shapes {
	newCoords := make(map[Coord]struct{})
	for c := range s.coords {
		newCoords[*c.Left()] = struct{}{}
		newCoords[*c.Right()] = struct{}{}
		newCoords[*c.Above()] = struct{}{}
		newCoords[*c.Below()] = struct{}{}
		newCoords[*c.Before()] = struct{}{}
		newCoords[*c.Behind()] = struct{}{}
	}
	for c := range s.coords {
		delete(newCoords, c)
	}

	result := NewShapes()
	numberOfNewShapes := len(newCoords)
	channel := make(chan *Shape, numberOfNewShapes)
	wg := sync.WaitGroup{}
	wg.Add(numberOfNewShapes)
	for c := range newCoords {
		go func(c Coord) {
			newShape := s.MustAddCube(&c)
			channel <- newShape.WithSmallestScore()
			wg.Done()
		}(c)
	}
	wg.Wait()
	close(channel)

	for shape := range channel {
		result.Add(shape)
	}

	return result
}

func (s *Shape) IsNeighbor(c *Coord) bool {
	for existing := range s.coords {
		for neighbor := range existing.Neighbors() {
			if c.Equals(neighbor) {
				return true
			}
		}
	}

	return false
}

func (s *Shape) Coords() []Coord {
	result := make([]Coord, 0, s.Size())
	for c := range s.coords {
		result = append(result, c)
	}

	return result
}

func (s *Shape) Transform(f func(Coord, Axis) (*Coord, error), axis Axis) (*Shape, error) {
	result := NewShape()
	result.coords = make(map[Coord]struct{}, s.Size())

	for c := range s.coords {
		newCoord, err := f(c, axis)
		if err != nil {
			return nil, err
		}
		result.coords[*newCoord] = struct{}{}
	}

	return result, nil
}

func (s *Shape) MustTransform(f func(Coord, Axis) (*Coord, error), axis Axis) *Shape {
	result, err := s.Transform(f, axis)
	if err != nil {
		panic(err)
	}

	return result
}

func (s *Shape) Rotate(axis Axis) (*Shape, error) {
	return s.Transform(
		func(c Coord, a Axis) (*Coord, error) {
			return c.Rotate(axis)
		},
		axis,
	)
}

func (s *Shape) MustRotate(axis Axis) *Shape {
	result, err := s.Rotate(axis)
	if err != nil {
		panic(err)
	}

	return result
}

func (s *Shape) Mirror(axis Axis) (*Shape, error) {
	return s.Transform(
		func(c Coord, a Axis) (*Coord, error) {
			return c.Mirror(axis)
		},
		axis,
	)
}

func (s *Shape) MustMirror(axis Axis) *Shape {
	result, err := s.Mirror(axis)
	if err != nil {
		panic(err)
	}

	return result
}

func (s *Shape) BoundingBox() (*Coord, *Coord) {
	var min, max Coord

	for _, c := range s.Coords() {
		for _, axis := range []Axis{XAxis, YAxis, ZAxis} {
			if c[axis] < min[axis] {
				min[axis] = c[axis]
			}
			if c[axis] > max[axis] {
				max[axis] = c[axis]
			}
		}
	}

	return &min, &max
}

func (s *Shape) AllPositiveCoords() *Shape {
	min, _ := s.BoundingBox()

	result := NewShape()
	result.coords = make(map[Coord]struct{}, s.Size())

	for c := range s.coords {
		newCoord := Coord{
			c[XAxis] - min[XAxis],
			c[YAxis] - min[YAxis],
			c[ZAxis] - min[ZAxis],
		}
		result.coords[newCoord] = struct{}{}
	}

	return result
}

func (s *Shape) Score() Score {
	result := make(Score, 0)
	var index uint64
	size := s.Size()
	sizeSquared := size * size

	for c := range s.AllPositiveCoords().coords {
		index = uint64(c[XAxis]) + uint64(c[YAxis]*size) + uint64(c[ZAxis]*sizeSquared)
		result[index] = true
	}

	return result
}

// true if s is equal to other
func (s *Shape) Equals(other *Shape) bool {
	return s.Cmp(other) == 0
}

// Compare left to right, return -1 if is left is smaller than right, 0 if left is equal to right and 1 if left is bigger than right
func (left *Shape) Cmp(right *Shape) int {
	return left.Score().Cmp(right.Score())
}

// returns the shape with the smallest score by rotating the original shape
func (s *Shape) WithSmallestScore() *Shape {
	// https://stackoverflow.com/questions/16452383/how-to-get-all-24-rotations-of-a-3-dimensional-array
	// RTTTRTTTRTTT
	// RTR
	// RTTTRTTTRTTT

	result := s
	turnedShape := s

	// RTTT
	turnedShape = turnedShape.MustRotate(XAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}

	// RTTT
	turnedShape = turnedShape.MustRotate(XAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}

	// RTTT
	turnedShape = turnedShape.MustRotate(XAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}

	// RTR
	turnedShape = turnedShape.MustRotate(XAxis).MustRotate(YAxis).MustRotate(XAxis)

	// RTTT
	turnedShape = turnedShape.MustRotate(XAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}

	// RTTT
	turnedShape = turnedShape.MustRotate(XAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}

	// RTTT
	turnedShape = turnedShape.MustRotate(XAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}
	turnedShape = turnedShape.MustRotate(YAxis)
	if result.Cmp(turnedShape) < 0 {
		result = turnedShape
	}

	return result.AllPositiveCoords()
}

// KeepGrowing returns all unique shapes starting from the initial Shape until the shapes reach the specified maxLen
func (initialShape *Shape) KeepGrowing(maxSize int, returnChannel chan Shapes) {
	if initialShape.Size() > maxSize {
		return
	}

	smallestScore := initialShape.WithSmallestScore()
	result := NewShapes()
	result.Add(smallestScore)

	grown := initialShape.Grow()

	requestChannel := make(chan Shapes, grown.NumberOfShapesWithSize(grown.maxSize))
	wg := sync.WaitGroup{}
	wg.Add(grown.NumberOfShapesWithSize(grown.maxSize))
	for _, shape := range grown.GetAllWithSize(grown.maxSize) {
		func(s *Shape) {
			defer wg.Done()
			s.KeepGrowing(maxSize, requestChannel)
		}(shape)
	}
	wg.Wait()
	close(requestChannel)

	for m := range requestChannel {
		result.Merge(&m)
	}

	returnChannel <- *result
}

func (s *Shape) String() string {
	coords := make([]string, 0)
	for c := range s.coords {
		coords = append(coords, c.String())
	}
	sort.Strings(coords)
	return strings.Join(coords, SEPARATOR)
}

func ShapeFromString(s string) (*Shape, error) {
	result := &Shape{coords: make(map[Coord]struct{})}
	coordStrings := strings.Split(s, SEPARATOR)
	for _, coordString := range coordStrings {
		coord, err := CoordFromString(coordString)
		if err != nil {
			return nil, err
		}
		result.coords[*coord] = struct{}{}
	}

	return result, nil
}
