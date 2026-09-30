package fixtures

import (
	"errors"
	"strings"
	"testing"
)

func TestParseCollectsAllErrors(t *testing.T) {
	src := "2024-03-01 Arsenal vs Chelsea\n" +
		"2024-03-02 Liverpool Everton\n" +
		"\n" +
		"03-08-2024 Chelsea vs Liverpool\n" +
		"2024-03-09 Everton vs \n"

	sched, err := Parse(strings.NewReader(src))
	if sched != nil {
		t.Errorf("expected nil schedule on error, got %+v", sched)
	}

	var errs ParseErrors
	if !errors.As(err, &errs) {
		t.Fatalf("expected ParseErrors, got %T: %v", err, err)
	}
	wantLines := []int{2, 4, 5}
	if len(errs) != len(wantLines) {
		t.Fatalf("got %d errors, want %d: %v", len(errs), len(wantLines), err)
	}
	for i, want := range wantLines {
		if errs[i].Line != want {
			t.Errorf("error %d: line %d, want %d", i, errs[i].Line, want)
		}
	}
}

func TestParseErrorsUnwrapToParseError(t *testing.T) {
	_, err := Parse(strings.NewReader("nonsense\n"))
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("errors.As did not find a *ParseError in %v", err)
	}
	if pe.Line != 1 || pe.Column != 1 {
		t.Errorf("got line %d column %d, want 1:1", pe.Line, pe.Column)
	}
}

func TestParseValid(t *testing.T) {
	src := "# opening weekend\n" +
		"2024-03-01 Arsenal vs Chelsea\n" +
		"\n" +
		"2024-03-08 Chelsea vs Liverpool\n"

	sched, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	if got := sched.Rounds(); got != 2 {
		t.Errorf("Rounds() = %d, want 2", got)
	}
	if got := sched.ByRound(2); len(got) != 1 || got[0].Home != "Chelsea" {
		t.Errorf("unexpected round 2: %+v", got)
	}
}

func TestParseErrorDoesNotShiftRounds(t *testing.T) {
	// The bad line in round 1 must not change how the following lines are
	// grouped; only the error list is affected.
	src := "2024-03-01 Arsenal vs Chelsea\n" +
		"bad\n" +
		"\n" +
		"2024-03-08 Chelsea Liverpool\n"

	_, err := Parse(strings.NewReader(src))
	var errs ParseErrors
	if !errors.As(err, &errs) || len(errs) != 2 {
		t.Fatalf("expected two errors, got %v", err)
	}
}

func TestParseErrorsMessageSeparatesEntries(t *testing.T) {
	_, err := Parse(strings.NewReader("x\ny\n"))
	if got := strings.Count(err.Error(), "\n\n"); got != 1 {
		t.Errorf("expected one blank-line separator, got %d in:\n%s", got, err)
	}
}
