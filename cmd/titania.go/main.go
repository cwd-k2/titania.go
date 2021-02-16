package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/cwd-k2/titania.go/pkg/tester"
)

const VERSION string = "v0.6.2"

func main() {
	// ターゲットのディレクトリと言語，quiet
	directories, languages, quiet := optparse()

	tester.SetQuiet(quiet)
	tester.SetLanguages(languages)

	outcomes := exec(directories)

	// 何もテストが実行されなかった場合
	if len(outcomes) == 0 {
		println("There's no test in (sub)directory[ies].")
		os.Exit(1)
	}

	final(outcomes)

	// buffering
	stdout := bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	// TODO: 全部メモリに持っておくのは辛いので (形式はそのままに) 分割して出力したい.
	if err := enc.Encode(outcomes); err != nil {
		panic(err)
	}
}
