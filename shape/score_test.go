package shape_test

import (
	"testing"

	. "github.com/munnik/cubes/shape"
)

func TestHash(t *testing.T) {
	s1 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{2, 0, 0})
	s2 := NewShape().MustAddCube(&Coord{1, 0, 0}).MustAddCube(&Coord{0, 1, 0})

	if s1.Score().Hash() == s2.Score().Hash() {
		t.Fatalf("Expected two different hashes but got %v and %v", s1.Score().Hash(), s2.Score().Hash())
	}
}
