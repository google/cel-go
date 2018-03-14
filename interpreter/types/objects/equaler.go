package objects

type Equaler interface {
	Equal(other interface{}) bool
}
