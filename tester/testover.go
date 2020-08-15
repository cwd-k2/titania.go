package tester

import "fmt"

func WrapUp(results []*TestInfo) {
}

func Detail(results []*TestInfo) {
	for _, info := range results {
		fmt.Println(info.UnitName)
		fmt.Println(info.CaseName)
		fmt.Println(info.Language)
		fmt.Println(info.Result)
		fmt.Println(info.Error)
		fmt.Println(info.Time)
	}
}
