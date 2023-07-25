package shape_test

import (
	"testing"

	. "github.com/munnik/cubes/shape"
)

func TestCmp(t *testing.T) {
	var f func() Shapes
	s1 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0}).WithSmallestScore()
	s2 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0}).WithSmallestScore()
	if s1.Cmp(s2) == 0 {
		t.Fatalf("Expected two different shapes but got %v and %v", s1, s2)
	}

	s1 = NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0}).MustAddCube(&Coord{0, 1, 0}).MustAddCube(&Coord{0, 0, 1}).WithSmallestScore()
	s2 = NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0}).MustAddCube(&Coord{0, -1, 0}).MustAddCube(&Coord{0, 0, 1}).WithSmallestScore()
	if s1.Cmp(s2) != 0 {
		t.Fatalf("Expected two equal shapes but got %v and %v", s1, s2)
	}
}

func TestLongestStraight(t *testing.T) {
	var f func() Shapes
	s1 := NewShape(f).MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0}).WithSmallestScore()
	if s1.LongestStraight() != 3 {
		t.Fatalf("Expected longest straight 3 but got %v", s1.LongestStraight())
	}
}
