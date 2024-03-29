package main

import (
	"fmt"
	"os"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

var (
	PAIZA_IO_URL     = "https://api.paiza.io"
	PAIZA_IO_API_KEY = "guest"
)

var opts struct {
	ShowDetail          bool
	InputFilePath       string
	Language            string
	StdoutFilePath      string
	StderrFilePath      string
	BuildStdoutFilePath string
	BuildStderrFilePath string
	ProgramFilePath     string
}

var (
	input       = os.Stdin
	stdout      = os.Stdout
	stderr      = os.Stderr
	buildstdout = os.Stderr
	buildstderr = os.Stderr
)

func init() {
	if val := os.Getenv("PAIZA_IO_URL"); val != "" {
		PAIZA_IO_URL = val
	}
	if val := os.Getenv("PAIZA_IO_API_KEY"); val != "" {
		PAIZA_IO_API_KEY = val
	}
}

func main() {
	optparse()

	if opts.Language == "" {
		opts.Language = runner.LangType(opts.ProgramFilePath)
	}

	run := runner.NewRunner(runner.Config{
		Host:   PAIZA_IO_URL,
		APIKey: PAIZA_IO_API_KEY,
	})

	sourcecode, _ := os.Open(opts.ProgramFilePath)
	defer sourcecode.Close()

	if opts.InputFilePath != "" {
		input, _ = os.Open(opts.InputFilePath)
		defer input.Close()
	}
	if opts.StdoutFilePath != "" {
		stdout, _ = os.Create(opts.StdoutFilePath)
		defer stdout.Close()
	}
	if opts.StderrFilePath != "" {
		stderr, _ = os.Create(opts.StderrFilePath)
		defer stderr.Close()
	}
	if opts.BuildStdoutFilePath != "" {
		buildstdout, _ = os.Create(opts.BuildStdoutFilePath)
		defer buildstdout.Close()
	}
	if opts.BuildStderrFilePath != "" {
		buildstderr, _ = os.Create(opts.BuildStderrFilePath)
		defer buildstderr.Close()
	}

	res, err := run.Run(&runner.OrderSpec{
		Language:    opts.Language,
		SourceCode:  sourcecode,
		Input:       input,
		Stdout:      stdout,
		Stderr:      stderr,
		BuildStdout: buildstdout,
		BuildStderr: buildstderr,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	if opts.ShowDetail {
		fmt.Println("====================")
		if res.BuildResult != "" {
			fmt.Printf("build_time:   %ss\n", res.BuildTime)
			fmt.Printf("build_memory: %dKB\n", res.BuildMemory/1024)
			fmt.Printf("build_result: %s\n", res.BuildResult)
		}
		fmt.Printf("time:         %ss\n", res.Time)
		fmt.Printf("memory:       %dKB\n", res.Memory/1024)
		fmt.Printf("result:       %s\n", res.Result)
	}

	if res.BuildExitCode != 0 {
		os.Exit(res.BuildExitCode)
	}
	os.Exit(res.ExitCode)
}
