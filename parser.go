package fixtures

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const separator = " vs "

// Parse reads a fixture schedule from r. It stops at the first malformed
// line and returns a *ParseError describing exactly where the problem is.
func Parse(r io.Reader) (*Schedule, error) {
	scanner := bufio.NewScanner(r)
	sched := &Schedule{}
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
			return nil, err
		}
		sched.Fixtures = append(sched.Fixtures, fx)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading schedule: %w", err)
	}
	return sched, nil
}

// parseLine parses a single non-blank, non-comment line in the form
// "<date> <home> vs <away>".
func parseLine(raw string, lineNo, round int) (Fixture, error) {
	idx := strings.IndexAny(raw, " \t")
	if idx == -1 {
		return Fixture{}, newParseError(lineNo, 1, raw,
			`expected a date followed by "<home> vs <away>", found only one field`)
	}

	dateStr := raw[:idx]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
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
