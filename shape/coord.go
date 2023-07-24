package shape

import (
	"fmt"
	"strconv"
	"strings"
)

type Coord [3]int

func (c Coord) Left() Coord {
	return Coord{c[XAxis] + 1, c[YAxis], c[ZAxis]}
}

func (c Coord) Right() Coord {
	return Coord{c[XAxis] - 1, c[YAxis], c[ZAxis]}
}

func (c Coord) Above() Coord {
	return Coord{c[XAxis], c[YAxis] + 1, c[ZAxis]}
}

func (c Coord) Below() Coord {
	return Coord{c[XAxis], c[YAxis] - 1, c[ZAxis]}
}

func (c Coord) Before() Coord {
	return Coord{c[XAxis], c[YAxis], c[ZAxis] + 1}
}

func (c Coord) Behind() Coord {
	return Coord{c[XAxis], c[YAxis], c[ZAxis] - 1}
}

// 90 degrees rotation around the specified axis
func (c *Coord) Rotate(axis Axis) (*Coord, error) {
	if axis == XAxis {
		return &Coord{c[XAxis], c[ZAxis], -c[YAxis]}, nil
	}
	if axis == YAxis {
		return &Coord{-c[ZAxis], c[YAxis], c[XAxis]}, nil
	}
	if axis == ZAxis {
		return &Coord{c[YAxis], -c[XAxis], c[ZAxis]}, nil
	}

	return nil, fmt.Errorf("unknown axis %d", axis)
}

func (c *Coord) MustRotate(axis Axis) *Coord {
	result, err := c.Rotate(axis)
	if err != nil {
		panic(err)
	}
	return result
}

// mirror using the plane orthogonal to the specified axis
func (c *Coord) Mirror(axis Axis) (*Coord, error) {
	if axis == XAxis {
		return &Coord{-c[XAxis], c[YAxis], c[ZAxis]}, nil
	}
	if axis == YAxis {
		return &Coord{c[XAxis], -c[YAxis], c[ZAxis]}, nil
	}
	if axis == ZAxis {
		return &Coord{c[XAxis], c[YAxis], -c[ZAxis]}, nil
	}

	return nil, fmt.Errorf("unknown axis %d", axis)
}

func (c *Coord) MustMirror(axis Axis) *Coord {
	result, err := c.Mirror(axis)
	if err != nil {
		panic(err)
	}
	return result
}

func (c *Coord) Neighbors() map[Coord]struct{} {
	result := make(map[Coord]struct{}, 6)
	result[c.Left()] = struct{}{}
	result[c.Right()] = struct{}{}
	result[c.Above()] = struct{}{}
	result[c.Below()] = struct{}{}
	result[c.Before()] = struct{}{}
	result[c.Behind()] = struct{}{}

	return result
}

func (c *Coord) String() string {
	return fmt.Sprintf("[%d %d %d]", c[XAxis], c[YAxis], c[ZAxis])
}

func (c *Coord) Equals(other *Coord) bool {
	return c[XAxis] == other[XAxis] && c[YAxis] == other[YAxis] && c[ZAxis] == other[ZAxis]
}

func (c *Coord) Subtract(other *Coord) *Coord {
	return &Coord{
		c[XAxis] - other[XAxis],
		c[YAxis] - other[YAxis],
		c[ZAxis] - other[ZAxis],
	}
}

func CoordFromString(s string) (*Coord, error) {
	var ok bool
	var err error
	if s, ok = strings.CutPrefix(s, "["); !ok {
		return nil, fmt.Errorf("coordinate should start with [")
	}
	if s, ok = strings.CutSuffix(s, "]"); !ok {
		return nil, fmt.Errorf("coordinate should end with ]")
	}
	fields := strings.Fields(s)
	result := &Coord{}
	if result[XAxis], err = strconv.Atoi(fields[XAxis]); err != nil {
		return nil, err
	}
	if result[YAxis], err = strconv.Atoi(fields[YAxis]); err != nil {
		return nil, err
	}
	if result[ZAxis], err = strconv.Atoi(fields[ZAxis]); err != nil {
		return nil, err
	}
	return result, nil
}
