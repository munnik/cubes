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

// Len is the number of cubes in the collection.
func (s *Shape) Len() int { return len(s.coords) }

func (s *Shape) AddCube(c *Coord) (*Shape, error) {
	if _, ok := s.coords[*c]; ok {
		return nil, fmt.Errorf("coord %v is already in this shape", c)
	}

	if !s.IsNeighbor(c) {
		return nil, fmt.Errorf("coord %v is not a neighbor of this shape", c)
	}

	result := NewShape()
	result.coords = make(map[Coord]struct{}, s.Len()+1)
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
func (s *Shape) Grow() Shapes {
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
	grownLength := s.Len() + 1
	numberOfNewShapes := len(newCoords)
	result[grownLength] = make(map[string]*Shape, numberOfNewShapes)
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
	result := make([]Coord, 0, s.Len())
	for c := range s.coords {
		result = append(result, c)
	}

	return result
}

func (s *Shape) Transform(f func(Coord, Axis) (*Coord, error), axis Axis) (*Shape, error) {
	result := NewShape()
	result.coords = make(map[Coord]struct{}, s.Len())

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
	result.coords = make(map[Coord]struct{}, s.Len())

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
	size := s.Len()
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
	newShape := s

	// RTTT
	newShape = newShape.MustRotate(XAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}

	// RTTT
	newShape = newShape.MustRotate(XAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}

	// RTTT
	newShape = newShape.MustRotate(XAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}

	// RTR
	newShape = newShape.MustRotate(XAxis)
	newShape = newShape.MustRotate(YAxis)
	newShape = newShape.MustRotate(XAxis)

	// RTTT
	newShape = newShape.MustRotate(XAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}

	// RTTT
	newShape = newShape.MustRotate(XAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}

	// RTTT
	newShape = newShape.MustRotate(XAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}
	newShape = newShape.MustRotate(YAxis)
	if result.Cmp(newShape) < 0 {
		result = newShape
	}

	return result.AllPositiveCoords()
}

// KeepGrowing returns all unique shapes starting from the initial Shape until the shapes reach the specified maxLen
func (initialShape *Shape) KeepGrowing(maxLen int, returnChannel chan Shapes) {
	if initialShape.Len() > maxLen {
		return
	}

	smallestScore := initialShape.WithSmallestScore()
	result := NewShapes()
	result.Add(smallestScore)

	grown := initialShape.Grow()

	requestChannel := make(chan Shapes, len(grown[initialShape.Len()+1]))
	wg := sync.WaitGroup{}
	wg.Add(len(grown[initialShape.Len()+1]))
	for _, shape := range grown[initialShape.Len()+1] {
		func(s *Shape) {
			defer wg.Done()
			s.KeepGrowing(maxLen, requestChannel)
		}(shape)
	}
	wg.Wait()
	close(requestChannel)

	for m := range requestChannel {
		result.Merge(m)
	}

	returnChannel <- result
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
