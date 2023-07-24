package store

import (
	"bufio"
	"fmt"
	"os"

	. "github.com/munnik/cubes/shape"
)

func WriteText(s *Shapes, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for size := 1; size <= s.MaxSize(); size++ {
		for _, shape := range s.AllWithSize(size) {
			fmt.Fprintln(f, shape)
		}
	}
}

func ReadText(path string) (*Shapes, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := NewShapes()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		shape, err := ShapeFromString(scanner.Text())
		if err != nil {
			return nil, err
		}
		result.Add(shape)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return result, nil
}
