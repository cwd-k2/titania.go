package tester

import (
	"github.com/cwd-k2/titania.go/pretty"
)

type View interface {
	Draw()
	Update(int)
}

type QuietView struct {
	name  string
	total int
	count int
}

type FancyView struct {
	name    string
	codes   int
	cases   int
	counts  []int
	indexes []string
}

func (testUnit *TestUnit) InitView(quiet bool) View {

	if quiet {
		view := new(QuietView)

		view.name = testUnit.Name
		view.total = len(testUnit.TestCodes) * len(testUnit.TestCases)

		return view

	} else {

		view := new(FancyView)

		view.name = testUnit.Name
		view.codes = len(testUnit.TestCodes)
		view.cases = len(testUnit.TestCases)

		view.counts = make([]int, len(testUnit.TestCodes))

		indexes := make([]string, 0, len(testUnit.TestCodes))
		for _, testCode := range testUnit.TestCodes {
			indexes = append(indexes, testCode.Name)
		}

		view.indexes = indexes

		return view
	}

}

func (view *FancyView) Draw() {

	pretty.Printf("%s\n", pretty.Bold(pretty.Cyan(view.name)))

	for _, index := range view.indexes {
		pretty.Printf(
			"[%s] %s %s\n",
			pretty.Yellow("WAIT"), "START",
			pretty.Bold(pretty.Blue(index)))
	}

}

func (view *FancyView) Update(position int) {
	view.counts[position]++

	pretty.Up(view.codes - position)
	pretty.Erase()

	if view.counts[position] == view.cases {
		pretty.Printf(
			"[%s] %02d/%02d %s",
			pretty.Green("DONE"),
			view.counts[position], view.cases,
			pretty.Bold(pretty.Blue(view.indexes[position])))
	} else {
		pretty.Printf(
			"[%s] %02d/%02d %s",
			pretty.Yellow("WAIT"),
			view.counts[position], view.cases,
			pretty.Bold(pretty.Blue(view.indexes[position])))
	}

	pretty.Down(view.codes - position)
	pretty.Beginning()
}

func (view *QuietView) Draw() {
	pretty.Printf(
		"[%s] %s\n",
		pretty.Green("LAUNCH"),
		pretty.Bold(pretty.Cyan(view.name)))
}

func (view *QuietView) Update(_ int) {
	view.count++

	if view.count == view.total {
		pretty.Printf(
			"[%s] %s\n",
			pretty.Yellow("FINISH"),
			pretty.Bold(pretty.Cyan(view.name)))
	}

}
