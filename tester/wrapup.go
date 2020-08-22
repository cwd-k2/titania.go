package tester

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cwd-k2/titania.go/pretty"
)

type Outcome struct {
	TestUnit string   `json:"test_unit"`
	Method   string   `json:"method"`
	Fruits   []*Fruit `json:"fruits"`
}

type Fruit struct {
	SourceCode string    `json:"source_code"`
	Language   string    `json:"language"`
	Details    []*Detail `json:"details"`
}

type Detail struct {
	TestCase string `json:"test_case"`
	Result   string `json:"result"`
	Time     string `json:"time"`
	Output   string `json:"output"`
	Error    string `json:"error"`
}

func Final(outcomes []*Outcome) {
	pretty.Printf("\n%s\n", pretty.Bold("ALL DONE"))

	for _, outcome := range outcomes {
		pretty.Printf("\n%s\n", pretty.Bold(pretty.Cyan(outcome.TestUnit)))

		for _, fruit := range outcome.Fruits {

			pretty.Printf("%s: %s\n", pretty.Bold(fruit.Language), pretty.Bold(pretty.Blue(fruit.SourceCode)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "PASS":
					pretty.Printf("%s: %s %ss\n", pretty.Green(detail.TestCase), pretty.Green(detail.Result), detail.Time)
				case "FAIL":
					pretty.Printf("%s: %s %ss\n", pretty.Yellow(detail.TestCase), pretty.Yellow(detail.Result), detail.Time)
				case "CLIENT ERROR":
					pretty.Printf("%s: %s\n", pretty.Magenta(detail.TestCase), pretty.Magenta(detail.Result))
				case "SERVER ERROR":
					pretty.Printf("%s: %s\n", pretty.Blue(detail.TestCase), pretty.Blue(detail.Result))
				case "TESTER ERROR":
					pretty.Printf("%s: %s\n", pretty.Bold(pretty.Red(detail.TestCase)), pretty.Bold(pretty.Red(detail.Result)))
				default:
					pretty.Printf("%s: %s\n", pretty.Red(detail.TestCase), pretty.Red(detail.Result))
				}
			}
		}
	}
}

func Print(outcomes []*Outcome) {
	// JSON 形式に変換
	rawout, err := json.MarshalIndent(outcomes, "", "  ")
	// JSON パース失敗
	if err != nil {
		panic(err)
	}

	// エスケープされた文字を戻す
	output, err := strconv.Unquote(strings.Replace(strconv.Quote(string(rawout)), `\\u`, `\u`, -1))
	// 変換失敗
	if err != nil {
		panic(err)
	}

	// 実行結果を JSON 形式で出力
	fmt.Println(string(output))
}
