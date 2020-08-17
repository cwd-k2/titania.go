package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cwd-k2/titania.go/pretty"
	"github.com/cwd-k2/titania.go/tester"
)

const VERSION = "0.0.0-alpha"

func main() {
	// ターゲットのディレクトリと言語
	directories, languages := OptParse()

	var outcomes []*tester.ShowUnit

	for _, dirname := range directories {
		testRoom := tester.NewTestRoom(dirname, languages)
		// 実行するテストがない
		if testRoom == nil {
			continue
		}

		fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Cyan(dirname)))
		fruits := testRoom.Exec()

		outcome := new(tester.ShowUnit)
		outcome.Name = dirname
		outcome.Fruits = fruits
		outcomes = append(outcomes, outcome)
	}

	if outcomes == nil {
		// 何もテストが実行されなかった場合
		println("Uh, OK, there's no test.")
	} else {
		fmt.Fprintf(os.Stderr, "\n%s\n", pretty.Bold("ALL DONE"))

		tester.WrapUp(outcomes)

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
		defer fmt.Println(string(output))
	}

}
