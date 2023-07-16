package store

import (
	"bufio"
	"fmt"
	"os"

	. "github.com/munnik/cubes/shape"
)

func WriteText(s Shapes, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for len := range s {
		for _, shape := range s[len] {
			fmt.Fprintln(f, shape)
		}
	}
}

func ReadText(path string) (Shapes, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(Shapes)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		shape, err := ShapeFromString(scanner.Text())
		if err != nil {
			return nil, err
		}
		if _, ok := result[shape.Len()]; !ok {
			result[shape.Len()] = make(map[string]*Shape)
		}
		result[shape.Len()][shape.Score().Hash()] = shape
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return result, nil
}
