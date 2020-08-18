package tester

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cwd-k2/titania.go/pretty"
)

type ShowUnit struct {
	Name   string      `json:"target"`
	Fruits []*ShowCode `json:"fruits"`
}

type ShowCode struct {
	Name     string      `json:"source_code"`
	Language string      `json:"language"`
	Details  []*ShowCase `json:"details"`
}

type ShowCase struct {
	Name   string `json:"test_case"`
	Result string `json:"result"`
	Time   string `json:"time"`
	OutPut string `json:"output"`
	Error  string `json:"error"`
}

// 流石に雑すぎる ちゃんと要約して
func WrapUp(outcomes []*ShowUnit) {
	pretty.Printf("\n%s\n", pretty.Bold("ALL DONE"))

	for _, outcome := range outcomes {
		pretty.Printf("\n%s\n", pretty.Bold(pretty.Cyan(outcome.Name)))

		for _, fruit := range outcome.Fruits {

			pretty.Printf("%s: %s\n", pretty.Bold(fruit.Language), pretty.Bold(pretty.Blue(fruit.Name)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "PASS":
					pretty.Printf("%s: %s %ss\n", pretty.Green(detail.Name), pretty.Green(detail.Result), detail.Time)
				case "FAIL":
					pretty.Printf("%s: %s %ss\n", pretty.Yellow(detail.Name), pretty.Yellow(detail.Result), detail.Time)
				case "CLIENT ERROR":
					pretty.Printf("%s: %s\n", pretty.Magenta(detail.Name), pretty.Magenta(detail.Result))
				case "SERVER ERROR":
					pretty.Printf("%s: %s\n", pretty.Blue(detail.Name), pretty.Blue(detail.Result))
				case "TESTER ERROR":
					pretty.Printf("%s: %s\n", pretty.Bold(pretty.Red(detail.Name)), pretty.Bold(pretty.Red(detail.Result)))
				default:
					pretty.Printf("%s: %s\n", pretty.Red(detail.Name), pretty.Red(detail.Result))
				}
			}
		}
	}
}

func Print(outcomes []*ShowUnit) {
	// JSON 形式に変換
	rawout, err := json.MarshalIndent(outcomes, "", "  ")
	// JSON パース失敗
	if err != nil {
		panic(err)
	}

	// エスケープされた文字を戻す
	output, err := strconv.Unquote(
		strings.Replace(strconv.Quote(string(rawout)), `\\u`, `\u`, -1))
	// 変換失敗
	if err != nil {
		panic(err)
	}

	// 実行結果を JSON 形式で出力
	fmt.Println(string(output))
}
