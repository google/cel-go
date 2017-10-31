package common

// Position is line/column position in an input source.
type Location struct {
	// Line is the line number of the position, starting with 1.
	line int

	// Column is the column number of the position, starting with 1.
	column int

	// Source is the optional name of the source input that was parsed. This is typically a file name.
	source Source
}

// NewLocation creates a new Location instance.
func NewLocation(s Source, l int, c int) Location {
	return Location{
		source: s,
		line:   l,
		column: c,
	}
}

func (l Location) Line() int {
	return l.line
}

func (l Location) Column() int {
	return l.column
}

func (l Location) Source() Source {
	return l.source
}

var NoLocation = Location{
	line:   0,
	column: 0,
	source: nil,
}
