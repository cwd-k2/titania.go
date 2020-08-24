package tester

import (
	"sync"
)

func Exec(directories, languages []string, async bool) []*Outcome {
	testMatters := MakeTestMatters(directories, languages)

	if len(testMatters) == 0 {
		return nil
	}

	outcomes := make([]*Outcome, len(testMatters))

	if async {

		wg := new(sync.WaitGroup)

		for i, testMatter := range testMatters {
			wg.Add(1)

			go func(i int, testMatter *TestMatter) {
				defer wg.Done()
				view := testMatter.InitView(true)
				outcome := testMatter.Exec(view)
				outcomes[i] = outcome
			}(i, testMatter)

		}

		wg.Wait()

	} else {

		for i, testMatter := range testMatters {
			view := testMatter.InitView(false)
			outcome := testMatter.Exec(view)
			outcomes[i] = outcome
		}

	}

	return outcomes
}
