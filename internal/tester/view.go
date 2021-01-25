package tester

import (
	"github.com/cwd-k2/titania.go/pkg/pretty"
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
	indices []string
}

func InitView(testUnit *TestUnit, quiet bool) View {

	if quiet {
		return &QuietView{testUnit.Name, len(testUnit.TestTargets) * len(testUnit.TestCases), 0}
	} else {

		indices := make([]string, len(testUnit.TestTargets))

		for i, testTarget := range testUnit.TestTargets {
			indices[i] = testTarget.Name
		}

		return &FancyView{
			testUnit.Name,
			len(testUnit.TestTargets),
			len(testUnit.TestCases),
			make([]int, len(testUnit.TestTargets)),
			indices,
		}
	}

}

func (view *FancyView) Draw() {

	pretty.Printf("%s\n", pretty.Bold(pretty.Cyan(view.name)))

	for _, index := range view.indices {
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
	pretty.Printf("%02d/%02d %s", view.counts[position], view.cases, pretty.Bold(pretty.Blue(view.indices[position])))

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
