package hikaku

type pass struct {
	level  int
	path   string
	index  int
	parent string
}

type valueDiff struct{}

type differenceContext []valueDiff

func newDifferenceContext() *differenceContext {
	return &differenceContext{}
}

type executionBuffer []func() error

func (e *executionBuffer) Add(cb func() error) {
	(*e) = append((*e), cb)
}

func (e *executionBuffer) Pop() func() error {
	var x func() error
	x, (*e) = (*e)[0], (*e)[1:]
	return x
}

func (e *executionBuffer) Len() int {
	return len(*e)
}

func newExecutionBuffer() *executionBuffer {
	return &executionBuffer{}
}
