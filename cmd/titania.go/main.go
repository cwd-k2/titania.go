package main

import (
	"os"
)

const VERSION string = "v0.9.0"

func main() {
	directories := optparse()

	uresults := exec(directories)

	// 何もテストが実行されなかった場合
	if len(uresults) == 0 {
		println("There's no test in (sub)directory[ies].")
		os.Exit(1)
	}

	final(uresults)
	printjson(uresults)
}
