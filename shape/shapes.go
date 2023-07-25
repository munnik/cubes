package shape

type Shapes interface {
	Len() int
	Add(shape Shape) Shapes
	Merge(other Shapes) Shapes
	GetAll() map[Score]*Shape
	GetAllWithSize(size ShapeSize) map[Score]*Shape
	MaxSize() ShapeSize
}
