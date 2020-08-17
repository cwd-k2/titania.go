package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

type View struct {
	units   int
	cases   int
	places  map[string]int
	counts  map[string]int
	indexes []string
}

func InitView(
	testCodes map[string]*TestCode,
	testCases map[string]*TestCase) *View {

	view := new(View)
	view.units = len(testCodes)
	view.cases = len(testCases)

	view.places = make(map[string]int)
	view.counts = make(map[string]int)

	i := 0
	var indexes []string
	for unitName := range testCodes {
		view.places[unitName] = i
		view.counts[unitName] = 0
		indexes = append(indexes, unitName)
		i++
	}

	view.indexes = indexes

	return view

}

func (view *View) Draw(unitName string) {
	fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Cyan(unitName)))

	for _, index := range view.indexes {
		fmt.Fprintf(
			os.Stderr, "[%s] %s %s\n",
			pretty.Yellow("WAIT"), "START",
			pretty.Bold(pretty.Blue(index)))
	}
}

func (view *View) Update(unitName string) {
	position := view.places[unitName]

	view.counts[unitName]++

	pretty.Up(view.units - position)
	pretty.Erase()
	if view.counts[unitName] == view.cases {
		fmt.Fprintf(
			os.Stderr, "[%s] %02d/%02d %s",
			pretty.Green("DONE"), view.counts[unitName], view.cases,
			pretty.Bold(pretty.Blue(unitName)))
	} else {
		fmt.Fprintf(
			os.Stderr, "[%s] %02d/%02d %s",
			pretty.Yellow("WAIT"), view.counts[unitName], view.cases,
			pretty.Bold(pretty.Blue(unitName)))
	}
	pretty.Down(view.units - position)
	pretty.Beginning()
}
