package main

import (
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func exec(directories []string) []*tester.Outcome {
	outcomes := make([]*tester.Outcome, 0)

	for _, dirname := range directories {
		// 設定
		tconf := tester.NewConfig(dirname)
		if tconf == nil {
			continue
		}

		tunit := tester.NewTestUnit(dirname, tconf)
		if tunit == nil {
			continue
		}

		outcomes = append(outcomes, tunit.Exec())
	}

	return outcomes
}
