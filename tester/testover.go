package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

// はよ作れ 流石に雑すぎる ちゃんと要約して
func WrapUp(dirname string, results []*TestInfo) {
	var before string

	fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Cyan(dirname)))
	// results はもうこの時点でソートされてる
	for _, info := range results {
		if info.UnitName != before {
			before = info.UnitName
			fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(pretty.Cyan(info.Language)), pretty.Bold(pretty.Blue(info.UnitName)))
		}
		switch info.Result {
		case "PASS":
			fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Green(info.CaseName), pretty.Green(info.Result), info.Time)
		case "WRONG ANSWER":
			fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Yellow(info.CaseName), pretty.Yellow(info.Result))
		case "CLIENT ERROR":
			fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Magenta(info.CaseName), pretty.Magenta(info.Result))
		default:
			fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Red(info.CaseName), pretty.Red(info.Result))
		}
	}
}
