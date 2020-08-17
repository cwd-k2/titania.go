package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

type View interface {
	Draw()
	Update(string)
}

type QuietView struct {
	name  string
	total int
	count int
}

type FancyView struct {
	name    string
	units   int
	cases   int
	places  map[string]int
	counts  map[string]int
	indexes []string
}

func InitView(
	name string,
	testCodes map[string]*TestCode,
	testCases map[string]*TestCase,
	quiet bool) View {

	if quiet {
		view := new(QuietView)

		view.name = name
		view.total = len(testCodes) * len(testCases)

		return view

	} else {

		view := new(FancyView)

		view.name = name
		view.units = len(testCodes)
		view.cases = len(testCases)

		view.places = make(map[string]int)
		view.counts = make(map[string]int)

		i := 0
		var indexes []string
		for codeName := range testCodes {
			view.places[codeName] = i
			view.counts[codeName] = 0
			indexes = append(indexes, codeName)
			i++
		}

		view.indexes = indexes

		return view
	}

}

func (view *FancyView) Draw() {

	fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Cyan(view.name)))

	for _, index := range view.indexes {
		fmt.Fprintf(
			os.Stderr, "[%s] %s %s\n",
			pretty.Yellow("WAIT"), "START",
			pretty.Bold(pretty.Blue(index)))
	}

}

func (view *FancyView) Update(codeName string) {
	position := view.places[codeName]
	view.counts[codeName]++

	pretty.Up(view.units - position)
	pretty.Erase()

	if view.counts[codeName] == view.cases {
		fmt.Fprintf(
			os.Stderr, "[%s] %02d/%02d %s",
			pretty.Green("DONE"), view.counts[codeName], view.cases,
			pretty.Bold(pretty.Blue(codeName)))
	} else {
		fmt.Fprintf(
			os.Stderr, "[%s] %02d/%02d %s",
			pretty.Yellow("WAIT"), view.counts[codeName], view.cases,
			pretty.Bold(pretty.Blue(codeName)))
	}

	pretty.Down(view.units - position)
	pretty.Beginning()
}

func (view *QuietView) Draw() {
	fmt.Fprintf(
		os.Stderr, "[%s] %s\n",
		pretty.Green("LAUNCH"),
		pretty.Bold(pretty.Cyan(view.name)))
}

func (view *QuietView) Update(codeName string) {
	view.count++

	if view.count == view.total {
		fmt.Fprintf(
			os.Stderr, "[%s] %s\n",
			pretty.Yellow("FINISH"),
			pretty.Bold(pretty.Cyan(view.name)))
	}

}
