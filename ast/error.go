package ast

type ErrorExpression struct {
	BaseExpression
}

func (e *ErrorExpression) String() string {
	return ToDebugString(e)
}

func (e *ErrorExpression) writeDebugString(w *debugWriter) {
	w.append("*!error!*")
	w.adorn(e)
}
