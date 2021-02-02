package main

import (
	"github.com/cwd-k2/titania.go/pkg/tester"
)

const VERSION string = "v0.5.1"

func main() {
	// ターゲットのディレクトリと言語，quiet
	directories, languages, quiet := optparse()

	tester.SetQuiet(quiet)
	tester.SetLanguages(languages)

	exec(directories)
}
