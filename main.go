package main

import (
	"fmt"
	"sync"

	. "github.com/munnik/cubes/shape"
	"github.com/munnik/cubes/store"
)

func main() {
	maxSize := 8

	var shapes Shapes
	// initialShapes = Shapes{}
	// initialShapes.Add(NewShape())
	var err error
	if shapes, err = store.ReadText("results/shapes.txt"); err != nil {
		shapes = Shapes{}
		shapes.Add(NewShape())
	}

	var maxLength int
	for len := range shapes {
		if len > maxLength {
			maxLength = len
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(len(shapes[maxLength]))
	c := make(chan Shapes, len(shapes[maxLength]))
	for _, shape := range shapes[maxLength] {
		go func(shape *Shape) {
			shape.KeepGrowing(maxSize, c)
			wg.Done()
		}(shape)
	}
	wg.Wait()
	close(c)
	for s := range c {
		shapes.Merge(s)
	}

	wg = sync.WaitGroup{}
	for len := range shapes {
		counter := 1
		for _, shape := range shapes[len] {
			wg.Add(1)
			go func(shape *Shape, len, counter int) {
				store.WriteImage(shape, 1024, 1024, fmt.Sprintf("results/shape_%02d_%05d.png", len, counter), 0.85)
				wg.Done()
			}(shape, len, counter)
			counter += 1
		}
	}
	wg.Wait()
	store.WriteText(shapes, "results/shapes.txt")
}
