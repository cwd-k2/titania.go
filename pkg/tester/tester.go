package tester

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[tester] ", log.Lshortfile)
