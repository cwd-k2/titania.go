package tester

import (
	"log"
	"os"
)

var (
	logger = log.New(os.Stderr, "[tester] ", log.Lshortfile)
	quiet  = false
)

func SetQuiet(b bool) {
	quiet = b
}
