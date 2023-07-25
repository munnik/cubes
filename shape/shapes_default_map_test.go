package shape_test

import (
	"testing"

	. "github.com/munnik/cubes/shape"
)

func TestAddDefaultMap(t *testing.T) {
	ShapesMap := NewShapesDefaultMap()

	var f func() Shapes
	s1 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})
	s3 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})

	ShapesMap.Add(*s1).Add(*s2).Add(*s3)

	if ShapesMap.Len() != 2 {
		t.Fatalf("Expected length to equal 2 but got %d", ShapesMap.Len())
	}
}

func TestMergeDefaultMap(t *testing.T) {
	ShapesMap1 := NewShapesDefaultMap()
	ShapesMap2 := NewShapesDefaultMap()

	var f func() Shapes
	s1 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})
	s3 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 0, 1})

	ShapesMap1.Add(*s1).Add(*s2)
	ShapesMap2.Add(*s2).Add(*s3)

	ShapesMap1.Merge(ShapesMap2)

	if ShapesMap1.Len() != 3 {
		t.Fatalf("Expected length to equal 3 but got %d", ShapesMap1.Len())
	}
}
