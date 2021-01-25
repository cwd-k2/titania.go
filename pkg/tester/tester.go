package tester

import (
	"sync"

	"github.com/cwd-k2/titania.go/pkg/viewer"
)

func Exec(directories, languages []string, async bool) []*Outcome {
	testUnits := MakeTestUnits(directories, languages)

	if len(testUnits) == 0 {
		return nil
	}

	outcomes := make([]*Outcome, len(testUnits))

	if async {

		wg := &sync.WaitGroup{}

		for i, testUnit := range testUnits {
			wg.Add(1)

			go func(i int, testUnit *TestUnit) {
				defer wg.Done()
				view := InitView(testUnit, true)
				outcome := testUnit.Exec(view)
				outcomes[i] = outcome
			}(i, testUnit)

		}

		wg.Wait()

	} else {

		for i, testUnit := range testUnits {
			view := InitView(testUnit, false)
			outcome := testUnit.Exec(view)
			outcomes[i] = outcome
		}

	}

	return outcomes
}

func InitView(testUnit *TestUnit, quiet bool) viewer.Viewer {

	if quiet {

		return viewer.NewQuietView(
			testUnit.Name,
			len(testUnit.TestTargets)*len(testUnit.TestCases),
		)

	} else {

		indices := make([]string, len(testUnit.TestTargets))

		for i, testTarget := range testUnit.TestTargets {
			indices[i] = testTarget.Name
		}

		return viewer.NewFancyView(
			testUnit.Name,
			len(testUnit.TestTargets),
			len(testUnit.TestCases),
			indices,
		)
	}

}
