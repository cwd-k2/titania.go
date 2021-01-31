package main

import (
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func exec(directories []string) []*tester.Outcome {
	outcomes := make([]*tester.Outcome, 0)

	for _, dirname := range directories {
		tunit := tester.NewTestUnit(dirname)
		if tunit == nil {
			continue
		}
		outcomes = append(outcomes, tunit.Exec())
	}

	return outcomes
}
