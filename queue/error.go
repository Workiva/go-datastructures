package queue

type DisposedError struct{}

func (de DisposedError) Error() string {
	return `Queue has been disposed.`
}
