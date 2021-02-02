package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func version() {
	fmt.Printf("titania.go %s\n", VERSION)
}

func usage() {
	print(`usage: titania.go [options] [directories]

directories:
  directories to test, that have titania.json.
  if not specified, titania.go will take all subdirectories as targets.

options:
  -h, --help                   show this help message
  -v, --version                show version
      --quiet                  quiet log
      --lang=lang1[,lang2,...] language[s] to test
`)
}

// パスを相対パスとして綺麗な形に
func cleanerpath(pwd, directory string) (string, error) {
	if filepath.IsAbs(directory) {
		return filepath.Rel(pwd, directory)
	} else {
		return filepath.Rel(pwd, filepath.Join(pwd, directory))
	}
}

// オプション解析
func optparse() ([]string, []string, bool) {
	var (
		args     = make([]string, 0)
		langs    = make([]string, 0)
		dirnames = make([]string, 0)
		quiet    = false
	)

	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			usage()
			os.Exit(0)
		} else if arg == "--version" || arg == "-v" {
			version()
			os.Exit(0)
		} else if arg == "--quiet" {
			quiet = true
		} else if strings.HasPrefix(arg, "--lang=") {
			langs = strings.Split(strings.Replace(arg, "--lang=", "", 1), ",")
		} else if strings.HasPrefix(arg, "-") {
			println("Unknown option:", arg)
			usage()
			os.Exit(1)
		} else {
			args = append(args, arg)
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		entries, err := ioutil.ReadDir(pwd)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				dirnames = append(dirnames, entry.Name())
			}
		}

	} else {

		for _, directory := range args {
			dirname, err := cleanerpath(pwd, directory)
			if err != nil {
				panic(err)
			}
			dirnames = append(dirnames, dirname)
		}

	}

	return dirnames, langs, quiet
}
