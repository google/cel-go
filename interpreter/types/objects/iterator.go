package objects

type Iterable interface {
	Iterator() Iterator
}

type Iterator interface {
	HasNext() bool
	Next() interface{}
}
