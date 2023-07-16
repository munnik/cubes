package store

import (
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
