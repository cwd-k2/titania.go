package main

import (
	"sync"

	"github.com/cwd-k2/titania.go/pkg/tester"
)

func Exec(directories, languages []string, async bool) []*tester.Outcome {
	testUnits := tester.MakeTestUnits(directories, languages, async)

	if len(testUnits) == 0 {
		return nil
	}

	outcomes := make([]*tester.Outcome, len(testUnits))

	if async {

		wg := &sync.WaitGroup{}

		for i, testUnit := range testUnits {
			wg.Add(1)

			go func(i int, testUnit *tester.TestUnit) {
				defer wg.Done()
				outcome := testUnit.Exec()
				outcomes[i] = outcome
			}(i, testUnit)

		}

		wg.Wait()

	} else {

		for i, testUnit := range testUnits {
			outcome := testUnit.Exec()
			outcomes[i] = outcome
		}

	}

	return outcomes
}
