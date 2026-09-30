# fixture-parser

A Go library for parsing plain-text sports fixture lists into structured
data.

Fixture lists for a local league, a five-a-side round robin, or a season
calendar usually start life as a text file someone edits by hand. Typos
creep in: a date in the wrong format, a missing "vs", a stray blank team
name. Most parsers report those errors as "invalid input" with no context.
This one reports the line, the column, and the exact spot to fix.

## Format

```
2024-03-01 Arsenal vs Chelsea
2024-03-02 Liverpool vs Everton

2024-03-08 Chelsea vs Liverpool
2024-03-09 Everton vs Arsenal
```

- One fixture per line: `YYYY-MM-DD Home vs Away`.
- A blank line starts a new round. Rounds are numbered from 1.
- Lines starting with `#` are comments and are ignored.

## Usage

```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/maryl544/fixture-parser"
)

func main() {
	src := `2024-03-01 Arsenal vs Chelsea
2024-03-02 Liverpool vs Everton

2024-03-08 Chelsea vs Liverpool`

	schedule, err := fixtures.Parse(strings.NewReader(src))
	if err != nil {
		log.Fatal(err)
	}

	for _, fx := range schedule.ByRound(1) {
		fmt.Printf("%s: %s vs %s\n", fx.Date.Format("Jan 2"), fx.Home, fx.Away)
	}
}
```

## Error messages

Given this input, where round 2 is missing its separator:

```
2024-03-01 Arsenal vs Chelsea

2024-03-08 Chelsea Liverpool
```

`Parse` returns:

```
line 3, column 12: expected " vs " separating home and away teams
2024-03-08 Chelsea Liverpool
           ^
```

The caret lands under the point where the parser expected to find the
separator, so the fix is obvious without re-reading the whole line.

## Status

Early skeleton. The parser handles the format above and reports every
malformed line in one pass, as a `ParseErrors` value. Use `errors.As` with
`*ParseError` to get at an individual problem. See the issue tracker for
what's planned next.

## License

MIT, see [LICENSE](LICENSE).
