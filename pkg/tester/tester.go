package tester

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logger    = log.New(os.Stderr, "[tester] ", log.Lshortfile|log.Ltime)
	quiet     = false
	languages = make([]string, 0)
	tmpdir    = filepath.Join(os.TempDir(), "titania.go", time.Now().Format("20060102150405"))
)

func SetTmpDir(dir string) {
	tmpdir = dir
}

// Set if the log's output should be quiet or not.
func SetQuiet(b bool) {
	quiet = b
}

// Set programming languages to test globally.
func SetLanguages(ls []string) {
	languages = append(languages, ls...)
}
