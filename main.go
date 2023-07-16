package main

import (
	"fmt"
	"sync"

	. "github.com/munnik/cubes/shape"
	"github.com/munnik/cubes/store"
)

func main() {
	maxSize := 8

	c := make(chan Shapes)
	go func() {
		NewShape().KeepGrowing(maxSize, c)
		close(c)
	}()
	allShapes := <-c

	wg := sync.WaitGroup{}
	for len, shapes := range allShapes {
		counter := 1
		for _, shape := range shapes {
			wg.Add(1)
			go func(shape *Shape, len, counter int) {
				store.WriteImage(shape, 1024, 1024, fmt.Sprintf("results/shape_%02d_%05d.png", len, counter), 0.85)
				wg.Done()
			}(shape, len, counter)
			counter += 1
		}
	}
	wg.Wait()
	store.WriteText(allShapes, "results/shapes.txt")
}
