package tester

import (
	"log"
	"os"
)

var (
	logger    = log.New(os.Stderr, "[tester] ", log.Lshortfile|log.Ltime)
	quiet     = false
	languages = make([]string, 0)
)

// Set if the log's output should be quiet or not.
func SetQuiet(b bool) {
	quiet = b
}

// Set programming languages to test globally.
func SetLanguages(ls []string) {
	languages = append(languages, ls...)
}
