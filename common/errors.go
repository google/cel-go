package common

import "fmt"

// Errors is the main error collector mechanism.
type Errors struct {
	errors []Error
}

// NewErrors returns a new Errors instance.
func NewErrors() *Errors {
	return &Errors{
		errors: []Error{},
	}
}

// ReportError captures an error report from the caller.
func (e *Errors) ReportError(l Location, format string, args ...interface{}) {
	e.reportErrorInstance(Error{
		Location: l,
		Message:  fmt.Sprintf(format, args...),
	})
}

// GetErrors returns all the errors that are accumulated so far.
func (e *Errors) GetErrors() []Error {
	return e.errors[:]
}

func (e *Errors) reportErrorInstance(err Error) {
	e.errors = append(e.errors, err)
}

func (e *Errors) String() string {
	result := ""
	for i, err := range e.errors {
		if i > 0 {
			result += "\n"
		}
		result += err.ToDisplayString()
	}
	return result
}
