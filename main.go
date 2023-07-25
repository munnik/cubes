package main

import (
	"flag"
	"fmt"
	"sync"

	. "github.com/munnik/cubes/shape"
	"github.com/munnik/cubes/store"
)

func main() {
	var maxSize int
	var fileName string
	var imagePath string
	var method string
	flag.IntVar(&maxSize, "n", 1, "Specify the maximum number of cubes a polycube can consist of. All unique polycubes from 1 to n cubes are calculated.")
	flag.StringVar(&fileName, "f", "", "File name to read existing polycubes from, new polycubes are written to this file. If no file name is specified no file is used to read from or write to.")
	flag.StringVar(&imagePath, "i", "", "Path were images should be written, existing images will be overwritten. If not specified no images will be generated")
	flag.StringVar(&method, "m", "DefaultMap", "Method to use to create a set of all the shapes created. Options are DefaultMap and LongestStraight.")
	flag.Parse()

	var NewShapes func() Shapes
	if method == "DefaultMap" {
		NewShapes = NewShapesDefaultMap
	} else if method == "LongestStraight" {
		NewShapes = NewShapesLongestStraightMap
	} else {
		panic("Unknown method specified")
	}

	var shapes Shapes
	shapes = NewShapes()
	var err error
	if shapes, err = store.ReadText(fileName, shapes); err != nil {
		shapes = NewShapes()
	}
	if shapes.Len() == 0 {
		shapes.Add(*NewShape(NewShapes))
	}

	currentMaxSize := shapes.MaxSize()

	wg := sync.WaitGroup{}
	currentShapes := shapes.GetAllWithSize(currentMaxSize)
	wg.Add(len(currentShapes))
	c := make(chan Shapes, len(currentShapes))
	for _, shape := range currentShapes {
		shape.SetNewShapesMethod(NewShapes)
		go func(shape *Shape) {
			shape.KeepGrowing(ShapeSize(maxSize), c)
			wg.Done()
		}(shape)
	}
	wg.Wait()
	close(c)
	for s := range c {
		shapes.Merge(s)
	}

	if fileName != "" {
		store.WriteText(shapes, fileName)
	}

	if imagePath != "" {
		wg = sync.WaitGroup{}
		for size := ShapeSize(1); size <= ShapeSize(maxSize); size++ {
			counter := 1
			for _, shape := range shapes.GetAllWithSize(size) {
				wg.Add(1)
				go func(shape *Shape, size ShapeSize, counter int) {
					store.WriteImage(shape, 1024, 1024, fmt.Sprintf("%s/shape_%02d_%015d.png", imagePath, size, counter), 0.85)
					wg.Done()
				}(shape, size, counter)
				counter += 1
			}
		}
		wg.Wait()
	}

	fmt.Printf("Found %d shapes with size %d\n", len(shapes.GetAllWithSize(ShapeSize(maxSize))), maxSize)
}
