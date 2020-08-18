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
	units   int
	cases   int
	counts  []int
	indexes []string
}

func InitView(name string, testCodes []*TestCode, testCases []*TestCase, quiet bool) View {

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

		view.counts = make([]int, len(testCodes))

		indexes := make([]string, 0, len(testCodes))
		for _, testCode := range testCodes {
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

	pretty.Up(view.units - position)
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

	pretty.Down(view.units - position)
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
