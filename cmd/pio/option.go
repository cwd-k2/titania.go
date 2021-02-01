package main

import (
	"os"
	"strings"
)

func usage() {
	msg := `Usage:
  pio [OPTIONS] PROGRAMFILE

Application Options:
      --stdin                read input from STDIN (overwritten by --input)
      --detail               show detail
      --input=FILE           read input from specified FILE
      --language=LANGUAGE    executed program's language (automatic detect)

Help Options:
  -h, --help                 Show this help message`
	println(msg)
}

func optparse() {
	var (
		used = -1
		args []string
	)

	for i, arg := range os.Args[1:] {
		if !strings.HasPrefix(arg, "-") {
			if i != used {
				args = append(args, arg)
			}
			continue
		}

		if arg == "--help" || arg == "-h" {
			usage()
			os.Exit(0)
		} else if arg == "--stdin" {
			opts.STDIN = true
		} else if arg == "--detail" {
			opts.Detail = true
		} else if strings.HasPrefix(arg, "--input=") {
			opts.InputFile = strings.ReplaceAll(arg, "--input=", "")
		} else if arg == "--input" {
			opts.InputFile = os.Args[i+2]
			used = i + 1
		} else if strings.HasPrefix(arg, "--language=") {
			opts.Language = strings.ReplaceAll(arg, "--language=", "")
		} else if arg == "--language" {
			opts.Language = os.Args[i+2]
			used = i + 1
		} else {
			println("No such option:", arg)
			usage()
			os.Exit(1)
		}
	}

	if len(args) != 1 {
		usage()
		os.Exit(1)
	} else {
		opts.Args.ProgramFile = args[0]
	}
}
