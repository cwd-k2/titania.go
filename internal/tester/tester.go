package tester

import (
	"sync"
)

func Exec(directories, languages []string, async bool) []*Outcome {
	testUnits := MakeTestUnits(directories, languages)

	if len(testUnits) == 0 {
		return nil
	}

	outcomes := make([]*Outcome, len(testUnits))

	if async {

		wg := new(sync.WaitGroup)

		for i, testUnit := range testUnits {
			wg.Add(1)

			go func(i int, tester *TestUnit) {
				defer wg.Done()
				view := InitView(tester, true)
				outcome := tester.Exec(view)
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
