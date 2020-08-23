package tester

import (
	"sync"
)

func Exec(directories, languages []string, async bool) []*Outcome {
	testTopics := MakeTestTopics(directories, languages)

	if len(testTopics) == 0 {
		return nil
	}

	outcomes := make([]*Outcome, len(testTopics))

	if async {

		wg := new(sync.WaitGroup)

		for i, testTopic := range testTopics {
			wg.Add(1)

			go func(i int, testTopic *TestTopic) {
				defer wg.Done()
				view := testTopic.InitView(true)
				outcome := testTopic.Exec(view)
				outcomes[i] = outcome
			}(i, testTopic)

		}

		wg.Wait()

	} else {

		for i, testTopic := range testTopics {
			view := testTopic.InitView(false)
			outcome := testTopic.Exec(view)
			outcomes[i] = outcome
		}

	}

	return outcomes
}
