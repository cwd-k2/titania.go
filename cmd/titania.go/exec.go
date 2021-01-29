package main

import (
	"sync"

	"github.com/cwd-k2/titania.go/pkg/tester"
)

func Exec(directories, languages []string, async bool) []*tester.Outcome {

	if async {
		var (
			tunits   = maketestunits(directories, languages)
			outcomes = make([]*tester.Outcome, len(tunits))
			wg       = &sync.WaitGroup{}
		)

		for i, tunit := range tunits {
			wg.Add(1)
			go func(i int, testUnit *tester.TestUnit) {
				defer wg.Done()
				outcome := testUnit.Exec()
				outcomes[i] = outcome
			}(i, tunit)
		}
		wg.Wait()

		return outcomes
	} else {
		outcomes := make([]*tester.Outcome, 0)

		for _, dirname := range directories {
			tunit := tester.NewTestUnit(dirname, languages, async)
			if tunit == nil {
				continue
			}
			outcomes = append(outcomes, tunit.Exec())
		}

		return outcomes
	}
}

func maketestunits(directories, languages []string) []*tester.TestUnit {
	tunits := make([]*tester.TestUnit, 0)
	for _, dirname := range directories {
		tunit := tester.NewTestUnit(dirname, languages, true)
		if tunit != nil {
			tunits = append(tunits, tunit)
		}
	}
	return tunits
}
