package store

import (
	"github.com/fogleman/ln/ln"
	. "github.com/munnik/cubes/shape"
)

func WriteImage(s *Shape, width, height float64, path string, cubeSize float64) {
	if cubeSize <= 0 || cubeSize > 1 {
		cubeSize = 1
	}

	scene := ln.Scene{}
	for _, c := range s.Coords() {
		scene.Add(
			ln.NewCube(
				ln.Vector{
					X: (float64)(c[XAxis]),
					Y: (float64)(c[YAxis]),
					Z: (float64)(c[ZAxis]),
				},
				ln.Vector{
					X: (float64)(c[XAxis]) + cubeSize,
					Y: (float64)(c[YAxis]) + cubeSize,
					Z: (float64)(c[ZAxis]) + cubeSize,
				},
			),
		)
	}

	// define camera parameters
	_, max := s.BoundingBox()
	max = max.Left().Above().Before().Left().Above().Before()
	eye := ln.Vector{
		X: (float64)(max[XAxis]) * 1.5,
		Y: (float64)(max[YAxis]) * 1.5,
		Z: (float64)(max[ZAxis]) * 1.5,
	} // camera position
	// eye := ln.Vector{X: 4, Y: 3, Z: 2} // camera position

	center := ln.Vector{X: 0, Y: 0, Z: 0} // camera looks at
	up := ln.Vector{X: 0, Y: 0, Z: 1}     // up direction

	// define rendering parameters
	fieldOfViewY := 80.0 // vertical field of view, degrees
	zNear := 1.0         // near z plane
	zFar := 100.0        // far z plane
	step := 0.01         // how finely to chop the paths for visibility testing

	// compute 2D paths that depict the 3D scene
	paths := scene.Render(eye, center, up, width, height, fieldOfViewY, zNear, zFar, step)

	// render the paths in an image
	paths.WriteToPNG(path, width, height)
}
