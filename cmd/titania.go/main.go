package main

import (
	"os"

	"github.com/cwd-k2/titania.go/internal/tester"
)

const VERSION = "v0.2.1"

func main() {
	// ターゲットのディレクトリと言語，async
	directories, languages, async := OptParse()
	outcomes := tester.Exec(directories, languages, async)

	// 何もテストが実行されなかった場合
	if outcomes == nil {
		println("Uh, OK, there's no test.")
		os.Exit(1)
	}

	tester.Final(outcomes)
	tester.Print(outcomes)
}
