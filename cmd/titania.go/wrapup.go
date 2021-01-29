package main

import (
	"github.com/cwd-k2/titania.go/internal/pkg/pretty"
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func final(outcomes []*tester.Outcome) {
	pretty.Printf("\n%s\n", pretty.Bold("ALL DONE"))

	for _, outcome := range outcomes {
		pretty.Printf("\n%s\n", pretty.Bold(pretty.Cyan(outcome.Name)))

		for _, fruit := range outcome.Fruits {
			pretty.Printf("%s: %s\n", pretty.Bold(fruit.Language), pretty.Bold(pretty.Blue(fruit.TestTarget)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "CLIENT ERROR":
					pretty.Printf("%s: %s\n", pretty.Magenta(detail.TestCase), pretty.Magenta(detail.Result))
				case "SERVER ERROR":
					pretty.Printf("%s: %s\n", pretty.Blue(detail.TestCase), pretty.Blue(detail.Result))
				case "TESTER ERROR":
					pretty.Printf("%s: %s\n", pretty.Bold(pretty.Red(detail.TestCase)), pretty.Bold(pretty.Red(detail.Result)))
				case "PASS":
					fallthrough
				case "FAIL":
					if detail.IsExpected {
						pretty.Printf("%s: %s %ss\n", pretty.Green(detail.TestCase), pretty.Green(detail.Result), detail.Time)
					} else {
						pretty.Printf("%s: %s %ss\n", pretty.Yellow(detail.TestCase), pretty.Yellow(detail.Result), detail.Time)
					}
				default:
					if detail.IsExpected {
						pretty.Printf("%s: %s\n", pretty.Green(detail.TestCase), pretty.Green(detail.Result))
					} else {
						pretty.Printf("%s: %s\n", pretty.Red(detail.TestCase), pretty.Red(detail.Result))
					}
				}
			}
		}
	}
}
