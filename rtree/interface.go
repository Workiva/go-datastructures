package rtree

type Rectangle interface {
	LowerLeft() (int32, int32)
	UpperRight() (int32, int32)
}
