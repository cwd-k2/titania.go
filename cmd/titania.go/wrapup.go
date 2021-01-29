package main

import (
	p "github.com/cwd-k2/titania.go/internal/pkg/pretty"
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func final(outcomes []*tester.Outcome) {
	p.Printf("\n%s\n", p.Bold("ALL DONE"))

	for _, outcome := range outcomes {
		p.Printf("\n%s\n", p.Bold(p.Cyan(outcome.Name)))

		for _, fruit := range outcome.Fruits {
			p.Printf("%s: %s\n", p.Bold(fruit.Language), p.Bold(p.Blue(fruit.TestTarget)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "CLIENT ERROR":
					p.Printf("%s: %s\n", p.Magenta(detail.TestCase), p.Magenta(detail.Result))
				case "SERVER ERROR":
					p.Printf("%s: %s\n", p.Blue(detail.TestCase), p.Blue(detail.Result))
				case "TESTER ERROR":
					p.Printf("%s: %s\n", p.Bold(p.Red(detail.TestCase)), p.Bold(p.Red(detail.Result)))
				case "PASS":
					fallthrough
				case "FAIL":
					if detail.IsExpected {
						p.Printf("%s: %s %ss\n", p.Green(detail.TestCase), p.Green(detail.Result), detail.Time)
					} else {
						p.Printf("%s: %s %ss\n", p.Yellow(detail.TestCase), p.Yellow(detail.Result), detail.Time)
					}
				default:
					if detail.IsExpected {
						p.Printf("%s: %s\n", p.Green(detail.TestCase), p.Green(detail.Result))
					} else {
						p.Printf("%s: %s\n", p.Red(detail.TestCase), p.Red(detail.Result))
					}
				}
			}
		}
	}
}
