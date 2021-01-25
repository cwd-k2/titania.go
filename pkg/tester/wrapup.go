package tester

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/pretty"
)

type Outcome struct {
	Name       string   `json:"name"`
	TestMethod string   `json:"test_method"`
	Fruits     []*Fruit `json:"fruits"`
}

type Fruit struct {
	TestTarget string    `json:"test_target"`
	Language   string    `json:"language"`
	Expect     string    `json:"expect"`
	Details    []*Detail `json:"details"`
}

type Detail struct {
	TestCase   string `json:"test_case"`
	Result     string `json:"result"`
	IsExpected bool   `json:"is_expected"`
	Time       string `json:"time"`
	Output     string `json:"output"`
	Error      string `json:"error"`
}

func Final(outcomes []*Outcome) {
	pretty.Printf("\n%s\n", pretty.Bold("ALL DONE"))

	for _, outcome := range outcomes {
		pretty.Printf("\n%s\n", pretty.Bold(pretty.Cyan(outcome.Name)))

		for _, fruit := range outcome.Fruits {
			pretty.Printf("%s: %s\n", pretty.Bold(fruit.Language), pretty.Bold(pretty.Blue(fruit.TestTarget)))

			for _, detail := range fruit.Details {
				switch detail.Result {
				case "CLIENT ERROR":
					pretty.Printf("%s: %s\n", pretty.Magenta(detail.TestCase), pretty.Magenta(detail.Result))
				case "SERVER ERROR":
					pretty.Printf("%s: %s\n", pretty.Blue(detail.TestCase), pretty.Blue(detail.Result))
				case "TESTER ERROR":
					pretty.Printf("%s: %s\n", pretty.Bold(pretty.Red(detail.TestCase)), pretty.Bold(pretty.Red(detail.Result)))
				case "PASS":
					fallthrough
				case "FAIL":
					if detail.IsExpected {
						pretty.Printf("%s: %s %ss\n", pretty.Green(detail.TestCase), pretty.Green(detail.Result), detail.Time)
					} else {
						pretty.Printf("%s: %s %ss\n", pretty.Yellow(detail.TestCase), pretty.Yellow(detail.Result), detail.Time)
					}
				default:
					if detail.IsExpected {
						pretty.Printf("%s: %s\n", pretty.Green(detail.TestCase), pretty.Green(detail.Result))
					} else {
						pretty.Printf("%s: %s\n", pretty.Red(detail.TestCase), pretty.Red(detail.Result))
					}
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
	fmt.Fprintln(os.Stdout, output)
}
