package shape

type Shapes interface {
	Len() int
	Add(shape Shape) Shapes
	GetAll() map[Score]*Shape
	GetAllWithSize(size ShapeSize) map[Score]*Shape
	Merge(other Shapes) Shapes
	MaxSize() ShapeSize
}
