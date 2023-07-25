package shape_test

import (
	"testing"

	. "github.com/munnik/cubes/shape"
)

func TestAdd(t *testing.T) {
	ShapesMap := NewShapesMap()

	s1 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})
	s3 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})

	ShapesMap.Add(*s1).Add(*s2).Add(*s3)

	if ShapesMap.Len() != 2 {
		t.Fatalf("Expected length to equal 2 but got %d", ShapesMap.Len())
	}
}

func TestMerge(t *testing.T) {
	ShapesMap1 := NewShapesMap()
	ShapesMap2 := NewShapesMap()

	s1 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})
	s3 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 0, 1})

	ShapesMap1.Add(*s1).Add(*s2)
	ShapesMap2.Add(*s2).Add(*s3)

	ShapesMap1.Merge(ShapesMap2)

	if ShapesMap1.Len() != 3 {
		t.Fatalf("Expected length to equal 3 but got %d", ShapesMap1.Len())
	}
}
