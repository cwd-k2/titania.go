package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

type ShowRoom struct {
	RoomName string      `json:"room_name"`
	Fruits   []*ShowUnit `json:"fruits"`
}

type ShowUnit struct {
	UnitName string      `json:"unit_name"`
	Language string      `json:"language"`
	Details  []*ShowCase `json:"details"`
}

type ShowCase struct {
	CaseName string `json:"case_name"`
	Result   string `json:"result"`
	Time     string `json:"time"`
	OutPut   string `json:"output"`
	Error    string `json:"error"`
}

// 流石に雑すぎる ちゃんと要約して
func WrapUp(outgo *ShowRoom) {
	fmt.Fprintf(os.Stderr, "\n%s\n", pretty.Bold(pretty.Cyan(outgo.RoomName)))
	for _, fruit := range outgo.Fruits {

		fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(fruit.Language), pretty.Bold(pretty.Blue(fruit.UnitName)))

		for _, detail := range fruit.Details {
			switch detail.Result {
			case "PASS":
				fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Green(detail.CaseName), pretty.Green(detail.Result), detail.Time)
			case "FAIL":
				fmt.Fprintf(os.Stderr, "%s: %s %ss\n", pretty.Yellow(detail.CaseName), pretty.Yellow(detail.Result), detail.Time)
			case "CLIENT ERROR":
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Magenta(detail.CaseName), pretty.Magenta(detail.Result))
			case "SERVER ERROR":
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Blue(detail.CaseName), pretty.Blue(detail.Result))
			case "TESTER ERROR":
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Bold(pretty.Red(detail.CaseName)), pretty.Bold(pretty.Red(detail.Result)))
			default:
				fmt.Fprintf(os.Stderr, "%s: %s\n", pretty.Red(detail.CaseName), pretty.Red(detail.Result))
			}
		}
	}
}
