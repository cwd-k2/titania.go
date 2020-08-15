package tester

import (
	"fmt"

	"github.com/cwd-k2/titania.go/pretty"
)

type TestView struct {
	Units     int
	Cases     int
	Total     int
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
	testView.Total = len(testUnits) * len(testCases)

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
		fmt.Printf("  [%s] %s\n", pretty.Blue("UNIT"), index)
		fmt.Printf("    [%s] %s\n", pretty.Yellow("WAIT"), "initializing...")
	}
}

func (testView *TestView) Refresh(testInfo *TestInfo) {
	unitName := testInfo.UnitName
	caseName := testInfo.CaseName
	position := testView.Positions[unitName]

	testView.Counters[unitName]++

	pretty.Up(2 * (testView.Units - position))
	pretty.Down(1)
	pretty.Erase()
	if testView.Counters[unitName] == testView.Cases {
		fmt.Printf("    [%s] %s\n", pretty.Green("DONE"), caseName)
	} else {
		fmt.Printf("    [%s] %s\n", pretty.Yellow("WAIT"), caseName)
	}
	pretty.Down(2 * (testView.Units - position - 1))
}
