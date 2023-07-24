package shape_test

import (
	"testing"

	. "github.com/munnik/cubes/shape"
)

func TestAdd(t *testing.T) {
	shapes := NewShapes()

	s1 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})
	s3 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})

	shapes.Add(s1).Add(s2).Add(s3)

	if shapes.Len() != 2 {
		t.Fatalf("Expected length to equal 2 but got %d", shapes.Len())
	}
}

func TestMerge(t *testing.T) {
	shapes1 := NewShapes()
	shapes2 := NewShapes()

	s1 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})
	s3 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 0, 1})

	shapes1.Add(s1).Add(s2)
	shapes2.Add(s2).Add(s3)

	shapes1.Merge(shapes2)

	if shapes1.Len() != 3 {
		t.Fatalf("Expected length to equal 3 but got %d", shapes1.Len())
	}
}
