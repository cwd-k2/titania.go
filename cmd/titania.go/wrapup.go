package main

import (
	. "github.com/cwd-k2/titania.go/internal/pkg/pretty"
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func final(outcomes []*tester.Outcome) {
	Printf("\n%s\n", Bold("ALL DONE"))

	for _, outcome := range outcomes {
		Printf("\n%s\n", Bold(Cyan(outcome.Name)))

		for _, fruit := range outcome.Fruits {
			Printf("%s: %s\n", Bold(fruit.Language), Bold(Blue(fruit.TestTarget)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "CLIENT ERROR":
					Printf("%s: %s\n", Magenta(detail.TestCase), Magenta(detail.Result))
				case "SERVER ERROR":
					Printf("%s: %s\n", Blue(detail.TestCase), Blue(detail.Result))
				case "TESTER ERROR":
					Printf("%s: %s\n", Bold(Red(detail.TestCase)), Bold(Red(detail.Result)))
				case "PASS":
					fallthrough
				case "FAIL":
					if detail.IsExpected {
						Printf("%s: %s %ss\n", Green(detail.TestCase), Green(detail.Result), detail.Time)
					} else {
						Printf("%s: %s %ss\n", Yellow(detail.TestCase), Yellow(detail.Result), detail.Time)
					}
				default:
					if detail.IsExpected {
						Printf("%s: %s\n", Green(detail.TestCase), Green(detail.Result))
					} else {
						Printf("%s: %s\n", Red(detail.TestCase), Red(detail.Result))
					}
				}
			}
		}
	}
}
