package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/tester"
)

var (
	prettyprint = false
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
      --pretty                 pretty print output json
      --quiet                  quiet log
      --lang=lang1[,lang2,...] language[s] to test
      --tmpdir=DIRNAME         set a directory where temporary files are put
      --maxjob=N               set a maximum number of jobs to run concurrently (N > 0)
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
func optparse() []string {
	var (
		args     = make([]string, 0)
		dirnames = make([]string, 0)
		used     = -1
	)

	pwd, _ := os.Getwd()

	for i, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			if i != used {
				args = append(args, arg)
			}
			continue
		}

		if strings.HasPrefix(arg, "--lang=") {
			langs := strings.Split(strings.TrimPrefix(arg, "--lang="), ",")
			tester.SetLanguages(langs)
			continue
		} else if strings.HasPrefix(arg, "--tmpdir=") {
			tmpdir, _ := cleanerpath(pwd, strings.TrimPrefix(arg, "--tmpdir="))
			if tmpdir != "" {
				tester.SetTmpDir(tmpdir)
			}
			continue
		} else if strings.HasPrefix(arg, "--maxjob=") {
			num, err := strconv.Atoi(strings.TrimPrefix(arg, "--maxjob="))
			if err != nil || num <= 0 {
				usage()
				os.Exit(1)
			}
			tester.SetMaxConcurrentJobs(num)
			continue
		}

		switch arg {
		case "--help", "-h":
			usage()
			os.Exit(0)
		case "--version", "-v":
			version()
			os.Exit(0)
		case "--quiet":
			tester.SetQuiet(true)
		case "--pretty":
			prettyprint = true
		case "--lang":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			langs := strings.Split(os.Args[i+2], ",")
			tester.SetLanguages(langs)
			used = i + 1
		case "--tmpdir":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			tmpdir, _ := cleanerpath(pwd, os.Args[i+2])
			if tmpdir != "" {
				tester.SetTmpDir(tmpdir)
			}
			used = i + 1
		case "--maxjob":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			num, err := strconv.Atoi(os.Args[i+2])
			if err != nil || num <= 0 {
				usage()
				os.Exit(1)
			}
			tester.SetMaxConcurrentJobs(num)
			used = i + 1
		default:
			println("No such option:", arg)
			usage()
			os.Exit(1)
		}
	}

	if len(args) == 0 {
		entries, _ := ioutil.ReadDir(pwd)
		for _, entry := range entries {
			if entry.IsDir() {
				dirnames = append(dirnames, entry.Name())
			}
		}
	} else {
		for _, directory := range args {
			dirname, _ := cleanerpath(pwd, directory)
			dirnames = append(dirnames, dirname)
		}
	}

	return dirnames
}
