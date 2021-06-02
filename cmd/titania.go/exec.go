package main

import (
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func exec(directories []string) []*tester.TestUnitResult {
	uresults := make([]*tester.TestUnitResult, 0)

	for _, dirname := range directories {
		// 設定
		tconf := tester.ReadConfig(dirname)
		if tconf == nil {
			continue
		}

		tunit := tester.ReadTestUnit(dirname, tconf)
		if tunit == nil {
			continue
		}

		uresults = append(uresults, tunit.Exec())
	}

	return uresults
}
