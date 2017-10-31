package operators

const (
	Conditional   = "_?_:_"
	LogicalAnd    = "_&&_"
	LogicalOr     = "_||_"
	LogicalNot    = "!_"
	In            = "_in_"
	Equals        = "_==_"
	NotEquals     = "_!=_"
	Less          = "_<_"
	LessEquals    = "_<=_"
	Greater       = "_>_"
	GreaterEquals = "_>=_"
	Add           = "_+_"
	Subtract      = "_-_"
	Multiply      = "_*_"
	Divide        = "_/_"
	Modulo        = "_%_"
	Negate        = "-_"
	Index         = "_[_]"
	Has           = "has"
	All           = "all"
	Exists        = "exists"
	ExistsOne     = "exists_one"
	Map           = "map"
	Filter        = "filter"
)

var operators = map[string]string{
	"+":  Add,
	"-":  Subtract,
	"*":  Multiply,
	"/":  Divide,
	"%":  Modulo,
	"in": In,
	"==": Equals,
	"!=": NotEquals,
	"<":  Less,
	"<=": LessEquals,
	">":  Greater,
	">=": GreaterEquals,
}

func Find(text string) (string, bool) {
	op, found := operators[text]
	return op, found
}
