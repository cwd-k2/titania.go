package main

import (
	"encoding/json"
	"os"

	"github.com/cwd-k2/titania.go/pkg/tester"
)

const VERSION = "v0.4.0"

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

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	if err := enc.Encode(outcomes); err != nil {
		panic(err)
	}
}
