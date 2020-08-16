package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/option"
	"github.com/cwd-k2/titania.go/pretty"
	"github.com/cwd-k2/titania.go/tester"
)

const VERSION = "0.0.0-alpha"

func main() {
	// ターゲットのディレクトリと言語
	directories, languages := option.Parse()

	// 実行されたターゲットの数
	i := 0
	details := make(map[string]interface{})

	for _, dirname := range directories {
		testRoom := tester.NewTestRoom(dirname, languages)
		// 実行するテストがない
		if testRoom == nil {
			continue
		}

		i++
		fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Green(dirname)))
		results := testRoom.Exec()

		// 要約 遅延実行でもええ
		defer tester.WrapUp(dirname, results)
		details[dirname] = results
	}

	if i == 0 {
		// 何もテストが実行されなかった場合
		println("Uh, OK, there's no test.")
	} else {
		output, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			panic(err)
		}
		// 実行結果を JSON 形式で出力
		fmt.Println(string(output))
	}

}
