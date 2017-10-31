package test

import (
	"fmt"
	"strings"
)

// Compare compares two strings, a for actual, e for expected, and returns true or false. The comparison is done,
// by filtering out whitespace (i.e. space, tabs and newline).
func Compare(a string, e string) bool {
	a = strings.Replace(a, " ", "", -1)
	a = strings.Replace(a, "\n", "", -1)
	a = strings.Replace(a, "\t", "", -1)

	e = strings.Replace(e, " ", "", -1)
	e = strings.Replace(e, "\n", "", -1)
	e = strings.Replace(e, "\t", "", -1)

	return a == e
}

// DiffMessage creates a diff dump message for test failures.
func DiffMessage(context string, actual interface{}, expected interface{}) string {
	result := fmt.Sprintf("FAILURE(%s)\n", context)
	result += "\n===== ACTUAL =====\n"
	result += strings.TrimSpace(fmt.Sprintf("%v", actual))
	result += "\n==== EXPECTED ====\n"
	result += strings.TrimSpace(fmt.Sprintf("%v", expected))
	result += "\n==================\n"
	return result
}
