package fixtures

import (
	"fmt"
	"strings"
)

// ParseError describes a problem found while parsing a schedule. It carries
// enough information to print the offending source line with a caret under
// the exact column, which is what makes tracking down a typo in a hand-edited
// fixture list fast instead of a guessing game.
type ParseError struct {
	Line    int    // 1-indexed line number in the source
	Column  int    // 1-indexed column where the problem starts
	Source  string // the raw text of the offending line
	Message string
}

func newParseError(line, column int, source, message string) *ParseError {
	return &ParseError{Line: line, Column: column, Source: source, Message: message}
}

// Error renders the problem as a message followed by the source line and a
// caret pointing at the column, e.g.:
//
//	line 4, column 12: expected " vs " separating home and away teams
//	2024-03-08 Chelsea Liverpool
//	           ^
func (e *ParseError) Error() string {
	caret := strings.Repeat(" ", e.Column-1) + "^"
	return fmt.Sprintf("line %d, column %d: %s\n%s\n%s", e.Line, e.Column, e.Message, e.Source, caret)
}
