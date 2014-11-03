package queue

type mockItem int

func (mi mockItem) Compare(other Item) bool {
	omi := other.(mockItem)
	return mi >= omi
}
