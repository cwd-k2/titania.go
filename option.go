package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/pretty"
)

func version() {
	pretty.Printf("titania.go %s\n", VERSION)
	os.Exit(1)
}

func usage() {
	pretty.Printf(`usage: titania.go [options] [directories]

targets:
  directories to test, that have titania.json.
  if not specified, titania.go will take all subdirectories as targets.

options:
  -h, --help                 show this help message
  -v, --version              show version
  --lang=lang1[,lang2[,...]] language[s] to test
`)
	os.Exit(1)
}

// パスを相対パスとして綺麗な形に
func cleanerPath(pwd, directory string) (string, error) {
	if filepath.IsAbs(directory) {
		return filepath.Rel(pwd, directory)
	} else {
		return filepath.Rel(pwd, filepath.Join(pwd, directory))
	}
}

// オプション解析
func OptParse() ([]string, []string, bool) {
	var args []string
	var async bool = false
	var languages []string
	var directories []string

	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" {
			usage()
		} else if arg == "--version" || arg == "-v" {
			version()
		} else if strings.HasPrefix(arg, "--lang=") {
			languages = strings.Split(strings.Replace(arg, "--lang=", "", 1), ",")
		} else if strings.HasPrefix(arg, "--async") {
			async = true
		} else if strings.HasPrefix(arg, "-") {
			pretty.Printf("Unknown option: %s\n", arg)
			usage()
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
				directories = append(directories, entry.Name())
			}
		}

	} else {

		for _, directory := range args {
			dirname, err := cleanerPath(pwd, directory)
			if err != nil {
				panic(err)
			}
			directories = append(directories, dirname)
		}

	}

	return directories, languages, async
}
