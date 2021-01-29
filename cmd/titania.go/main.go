package main

import (
	"os"
)

const VERSION = "v0.3.0-alpha"

func main() {
	// ターゲットのディレクトリと言語，async
	dirnames, languages, async := OptParse()
	outcomes := Exec(dirnames, languages, async)

	// 何もテストが実行されなかった場合
	if outcomes == nil {
		println("Uh, OK, there's no test.")
		os.Exit(1)
	}

	Final(outcomes)
	Print(outcomes)
}
