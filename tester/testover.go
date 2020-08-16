package tester

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
)

// はよ作れ
func WrapUp(dirname string, results []*TestInfo) {
	fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Green(dirname)))
}
