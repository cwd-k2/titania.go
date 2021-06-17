package main

import (
	"os"

	"github.com/cwd-k2/titania.go/pkg/tester"
)

const VERSION string = "v0.8.0"

func main() {
	directories := optparse()

	tester.SetQuiet(quiet)
	tester.SetLanguages(langs)
	if tmpdir != "" {
		tester.SetTmpDir(tmpdir)
	}

	uresults := exec(directories)

	// 何もテストが実行されなかった場合
	if len(uresults) == 0 {
		println("There's no test in (sub)directory[ies].")
		os.Exit(1)
	}

	final(uresults)
	printjson(uresults)
}
