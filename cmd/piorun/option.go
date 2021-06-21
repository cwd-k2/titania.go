package main

import (
	"os"
	"strings"
)

func usage() {
	println(`Usage: piorun [OPTIONS] PROGRAMFILE

Options:
      --detail               show detail
      --input=FILE           read input from specified FILE
      --stdout=FILE          write stdout to specified FILE
      --stderr=FILE          write stderr to specified FILE
      --build-stdout=FILE    write build stdout to specified FILE
      --build-stderr=FILE    write build stderr to specified FILE
      --language=LANGUAGE    override program's language
  -h, --help                 Show this help message`)
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

		if strings.HasPrefix(arg, "--input=") {
			opts.InputFilePath = strings.TrimPrefix(arg, "--input=")
			continue
		} else if strings.HasPrefix(arg, "--stdout=") {
			opts.StdoutFilePath = strings.TrimPrefix(arg, "--stdout=")
			continue
		} else if strings.HasPrefix(arg, "--stderr=") {
			opts.StderrFilePath = strings.TrimPrefix(arg, "--stderr=")
			continue
		} else if strings.HasPrefix(arg, "--build-stdout=") {
			opts.BuildStdoutFilePath = strings.TrimPrefix(arg, "--build-stdout=")
			continue
		} else if strings.HasPrefix(arg, "--build-stderr=") {
			opts.BuildStderrFilePath = strings.TrimPrefix(arg, "--build-stderr=")
			continue
		} else if strings.HasPrefix(arg, "--language=") {
			opts.Language = strings.TrimPrefix(arg, "--language=")
			continue
		}

		switch arg {
		case "--help", "-h":
			usage()
			os.Exit(0)
		case "--detail":
			opts.ShowDetail = true
		case "--input":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			opts.InputFilePath = os.Args[i+2]
			used = i + 1
		case "--stdout":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			opts.StdoutFilePath = os.Args[i+2]
			used = i + 1
		case "--stderr":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			opts.StderrFilePath = os.Args[i+2]
			used = i + 1
		case "--build-stdout":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			opts.BuildStdoutFilePath = os.Args[i+2]
			used = i + 1
		case "--build-stderr":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			opts.BuildStderrFilePath = os.Args[i+2]
			used = i + 1
		case "--language":
			if len(os.Args) < i+3 {
				usage()
				os.Exit(1)
			}
			opts.Language = os.Args[i+2]
			used = i + 1
		default:
			println("No such option:", arg)
			usage()
			os.Exit(1)
		}
	}

	if len(args) != 1 {
		usage()
		os.Exit(1)
	} else {
		opts.ProgramFilePath = args[0]
	}
}
