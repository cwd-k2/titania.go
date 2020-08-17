package main

import (
	"os"

	"github.com/cwd-k2/titania.go/tester"
)

const VERSION = "0.0.0-alpha"

func main() {
	// ターゲットのディレクトリと言語
	directories, languages := OptParse()
	outcomes := tester.Execute(directories, languages)

	// 何もテストが実行されなかった場合
	if outcomes == nil {
		println("Uh, OK, there's no test.")
		os.Exit(1)
	}

	tester.WrapUp(outcomes)
	defer tester.OutPut(outcomes)
}
