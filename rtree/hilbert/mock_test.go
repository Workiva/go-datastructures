package hilbert

type mockRectangle struct {
	xlow, ylow, xhigh, yhigh int32
}

func (mr *mockRectangle) LowerLeft() (int32, int32) {
	return mr.xlow, mr.ylow
}

func (mr *mockRectangle) UpperRight() (int32, int32) {
	return mr.xhigh, mr.yhigh
}

func newMockRectangle(xlow, ylow, xhigh, yhigh int32) *mockRectangle {
	return &mockRectangle{
		xlow:  xlow,
		ylow:  ylow,
		xhigh: xhigh,
		yhigh: yhigh,
	}
}
