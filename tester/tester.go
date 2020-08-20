package tester

import (
	"sync"
)

func Exec(directories, languages []string, async bool) []*Outcome {
	targets := MakeTargets(directories, languages)

	if len(targets) == 0 {
		return nil
	}

	outcomes := make([]*Outcome, len(targets))

	if async {

		wg := new(sync.WaitGroup)

		for i, target := range targets {
			wg.Add(1)

			go func(i int, target *Target) {
				defer wg.Done()
				view := target.InitView(true)
				outcome := target.Exec(view)
				outcomes[i] = outcome
			}(i, target)

		}

		wg.Wait()

	} else {

		for i, target := range targets {
			view := target.InitView(false)
			outcome := target.Exec(view)
			outcomes[i] = outcome
		}

	}

	return outcomes
}
