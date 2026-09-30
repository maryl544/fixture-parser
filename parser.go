package fixtures

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const separator = " vs "

// Parse reads a fixture schedule from r. It keeps going after a malformed
// line so that every problem in a hand-edited file can be fixed in one pass.
// If any line is malformed it returns a nil schedule and a ParseErrors value
// with one *ParseError per bad line.
func Parse(r io.Reader) (*Schedule, error) {
	scanner := bufio.NewScanner(r)
	sched := &Schedule{}
	var errs ParseErrors
	round := 1
	sawBlank := false

	for lineNo := 1; scanner.Scan(); lineNo++ {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			sawBlank = true
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		if sawBlank && len(sched.Fixtures) > 0 {
			round++
		}
		sawBlank = false

		fx, err := parseLine(raw, lineNo, round)
		if err != nil {
			// A bad line still counts toward round inference, so the
			// lines after it are not reported against the wrong round.
			errs = append(errs, err)
			continue
		}
		sched.Fixtures = append(sched.Fixtures, fx)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading schedule: %w", err)
	}
	if len(errs) > 0 {
		return nil, errs
	}
	return sched, nil
}

// parseLine parses a single non-blank, non-comment line in the form
// "<date> <home> vs <away>".
func parseLine(raw string, lineNo, round int) (Fixture, *ParseError) {
	idx := strings.IndexAny(raw, " \t")
	if idx == -1 {
		return Fixture{}, newParseError(lineNo, 1, raw,
			`expected a date followed by "<home> vs <away>", found only one field`)
	}

	dateStr := raw[:idx]
	date, terr := time.Parse("2006-01-02", dateStr)
	if terr != nil {
		return Fixture{}, newParseError(lineNo, 1, raw,
			fmt.Sprintf("invalid date %q: expected format YYYY-MM-DD", dateStr))
	}

	rest := strings.TrimLeft(raw[idx:], " \t")
	restCol := len(raw) - len(rest) // 0-indexed offset where rest begins

	sepIdx := strings.Index(rest, separator)
	if sepIdx == -1 {
		return Fixture{}, newParseError(lineNo, restCol+1, raw,
			`expected " vs " separating home and away teams`)
	}

	home := strings.TrimSpace(rest[:sepIdx])
	if home == "" {
		return Fixture{}, newParseError(lineNo, restCol+1, raw, "home team name is empty")
	}

	away := strings.TrimSpace(rest[sepIdx+len(separator):])
	if away == "" {
		awayCol := restCol + sepIdx + len(separator) + 1
		return Fixture{}, newParseError(lineNo, awayCol, raw, "away team name is empty")
	}

	return Fixture{Date: date, Home: home, Away: away, Round: round}, nil
}
