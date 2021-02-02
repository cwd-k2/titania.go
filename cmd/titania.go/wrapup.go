package main

import (
	"io"

	. "github.com/cwd-k2/titania.go/internal/pkg/pretty"
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func printoutcome(w io.Writer, outcome *tester.Outcome) {
	Fprintf(w, "\n%s\n", Bold(Cyan(outcome.Name)))

	for _, fruit := range outcome.Fruits {
		printfruit(w, fruit)
	}
}

func printfruit(w io.Writer, fruit *tester.Fruit) {
	Fprintf(w, "%s: %s\n", Bold(fruit.Language), Bold(Blue(fruit.TestTarget)))

	for _, detail := range fruit.Details {
		printdetail(w, detail)
	}
}

func printdetail(w io.Writer, detail *tester.Detail) {
	switch detail.Result {
	case "CLIENT ERROR":
		Fprintf(w, "%s: %s\n", Magenta(detail.TestCase), Magenta(detail.Result))
	case "SERVER ERROR":
		Fprintf(w, "%s: %s\n", Blue(detail.TestCase), Blue(detail.Result))
	case "TESTER ERROR":
		Fprintf(w, "%s: %s\n", Bold(Red(detail.TestCase)), Bold(Red(detail.Result)))
	case "PASS":
		fallthrough
	case "FAIL":
		if detail.IsExpected {
			Fprintf(w, "%s: %s %ss\n", Green(detail.TestCase), Green(detail.Result), detail.Time)
		} else {
			Fprintf(w, "%s: %s %ss\n", Yellow(detail.TestCase), Yellow(detail.Result), detail.Time)
		}
	default:
		if detail.IsExpected {
			Fprintf(w, "%s: %s\n", Green(detail.TestCase), Green(detail.Result))
		} else {
			Fprintf(w, "%s: %s\n", Red(detail.TestCase), Red(detail.Result))
		}
	}
}
