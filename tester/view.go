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

func (testTopic *TestTopic) InitView(quiet bool) View {

	if quiet {
		view := new(QuietView)

		view.name = testTopic.Name
		view.total = len(testTopic.TestTargets) * len(testTopic.TestCases)

		return view

	} else {

		view := new(FancyView)

		view.name = testTopic.Name
		view.codes = len(testTopic.TestTargets)
		view.cases = len(testTopic.TestCases)

		view.counts = make([]int, len(testTopic.TestTargets))

		indexes := make([]string, 0, len(testTopic.TestTargets))
		for _, testTarget := range testTopic.TestTargets {
			indexes = append(indexes, testTarget.Name)
		}

		view.indexes = indexes

		return view
	}

}

func (view *FancyView) Draw() {

	pretty.Printf("%s\n", pretty.Bold(pretty.Cyan(view.name)))

	for _, index := range view.indexes {
		pretty.Printf("[%s] %s %s\n", pretty.Yellow("WAIT"), "START", pretty.Bold(pretty.Blue(index)))
	}

}

func (view *FancyView) Update(position int) {
	view.counts[position]++

	pretty.Up(view.codes - position)
	pretty.Erase()

	if view.counts[position] == view.cases {
		pretty.Printf("[%s] ", pretty.Green("DONE"))
	} else {
		pretty.Printf("[%s] ", pretty.Yellow("WAIT"))
	}
	pretty.Printf("%02d/%02d %s", view.counts[position], view.cases, pretty.Bold(pretty.Blue(view.indexes[position])))

	pretty.Down(view.codes - position)
	pretty.Beginning()
}

func (view *QuietView) Draw() {
	pretty.Printf("[%s] %s\n", pretty.Green("LAUNCH"), pretty.Bold(pretty.Cyan(view.name)))
}

func (view *QuietView) Update(_ int) {
	view.count++

	if view.count == view.total {
		pretty.Printf("[%s] %s\n", pretty.Yellow("FINISH"), pretty.Bold(pretty.Cyan(view.name)))
	}

}
