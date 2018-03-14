package objects

type Indexer interface {
	Get(index interface{}) (interface{}, error)
}
