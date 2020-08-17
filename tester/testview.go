package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

type TestView struct {
	Units     int
	Cases     int
	Positions map[string]int
	Counters  map[string]int
	Indexes   []string
}

func InitTestView(
	testUnits map[string]*TestUnit,
	testCases map[string]*TestCase) *TestView {

	testView := new(TestView)
	testView.Units = len(testUnits)
	testView.Cases = len(testCases)

	testView.Positions = make(map[string]int)
	testView.Counters = make(map[string]int)

	i := 0
	var indexes []string
	for unitName := range testUnits {
		testView.Positions[unitName] = i
		testView.Counters[unitName] = 0
		indexes = append(indexes, unitName)
		i++
	}

	testView.Indexes = indexes

	return testView

}

func (testView *TestView) Start() {
	for _, index := range testView.Indexes {
		fmt.Fprintf(
			os.Stderr,
			"[%s] %s %s\n",
			pretty.Yellow("WAIT"),
			"START",
			pretty.Bold(pretty.Blue(index)))
	}
}

func (testView *TestView) Refresh(unitName string) {
	position := testView.Positions[unitName]

	testView.Counters[unitName]++

	pretty.Up(testView.Units - position)
	pretty.Erase()
	if testView.Counters[unitName] == testView.Cases {
		fmt.Fprintf(
			os.Stderr,
			"[%s] %02d/%02d %s",
			pretty.Green("DONE"),
			testView.Counters[unitName],
			testView.Cases,
			pretty.Bold(pretty.Blue(unitName)))
	} else {
		fmt.Fprintf(
			os.Stderr,
			"[%s] %02d/%02d %s",
			pretty.Yellow("WAIT"),
			testView.Counters[unitName],
			testView.Cases,
			pretty.Bold(pretty.Blue(unitName)))
	}
	pretty.Down(testView.Units - position)
	pretty.Beginning()
}
