package main

import (
	"encoding/json"
	"os"
)

const VERSION = "v0.3.0"

func main() {
	// ターゲットのディレクトリと言語，async
	dirnames, languages, async := optparse()

	outcomes := exec(dirnames, languages, async)

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
