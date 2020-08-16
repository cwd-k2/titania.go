package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

type TestOver struct {
	UnitName string      `json:"unit_name"`
	Language string      `json:"language"`
	Details  []*TestInfo `json:"details"`
}

// 流石に雑すぎる ちゃんと要約して
func WrapUp(dirname string, results []*TestOver) {
	fmt.Fprintf(os.Stderr, "\n%s\n", pretty.Bold(pretty.Cyan(dirname)))
	for _, over := range results {
		fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(over.Language), pretty.Bold(pretty.Blue(over.UnitName)))
		for _, info := range over.Details {
			switch info.Result {
			case "PASS":
				fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Green(info.CaseName), pretty.Green(info.Result), info.Time)
			case "FAIL":
				fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Yellow(info.CaseName), pretty.Yellow(info.Result), info.Time)
			case "CLIENT ERROR":
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Magenta(info.CaseName), pretty.Magenta(info.Result))
			case "SERVER ERROR":
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Blue(info.CaseName), pretty.Blue(info.Result))
			case "TESTER ERROR":
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(pretty.Red(info.CaseName)), pretty.Bold(pretty.Red(info.Result)))
			default:
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Red(info.CaseName), pretty.Red(info.Result))
			}
		}
	}
}
