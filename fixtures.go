// Package fixtures parses plain-text sports fixture lists into structured
// data.
//
// The source format is deliberately simple, because it's meant to be typed
// and edited by hand:
//
//	2024-03-01 Arsenal vs Chelsea
//	2024-03-02 Liverpool vs Everton
//
//	2024-03-08 Chelsea vs Liverpool
//	2024-03-09 Everton vs Arsenal
//
// A blank line starts a new round. Lines beginning with # are comments and
// are skipped.
package fixtures

import "time"

// Fixture is a single scheduled match between two teams.
type Fixture struct {
	Date  time.Time
	Home  string
	Away  string
	Round int
}

// Schedule is an ordered list of fixtures, grouped implicitly by Round.
type Schedule struct {
	Fixtures []Fixture
}

// ByRound returns the fixtures belonging to the given round, in the order
// they appeared in the source. Rounds are numbered from 1.
func (s *Schedule) ByRound(round int) []Fixture {
	var out []Fixture
	for _, fx := range s.Fixtures {
		if fx.Round == round {
			out = append(out, fx)
		}
	}
	return out
}

// Rounds returns the highest round number present in the schedule, or 0 if
// the schedule has no fixtures.
func (s *Schedule) Rounds() int {
	max := 0
	for _, fx := range s.Fixtures {
		if fx.Round > max {
			max = fx.Round
		}
	}
	return max
}
