package tester

import (
	"log"
	"os"
)

var (
	logger    = log.New(os.Stderr, "[tester] ", log.Lshortfile)
	quiet     = false
	languages = make([]string, 0)
)

func SetQuiet(b bool) {
	quiet = b
}

func SetLanguages(ls []string) {
	languages = append(languages, ls...)
}
