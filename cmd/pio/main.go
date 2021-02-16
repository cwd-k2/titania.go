package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cwd-k2/titania.go/internal/pkg/langtype"
	"github.com/cwd-k2/titania.go/pkg/paizaio"
)

var (
	PAIZA_IO_URL     = "https://api.paiza.io"
	PAIZA_IO_API_KEY = "guest"
)

var opts struct {
	STDIN     bool
	Detail    bool
	InputFile string
	Language  string
	Args      struct {
		ProgramFile string
	}
}

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

	filename := opts.Args.ProgramFile

	sourceBS, err := ioutil.ReadFile(filename)
	if err != nil {
		println(err.Error())
		return
	}
	source := string(sourceBS)

	var language string
	if opts.Language != "" {
		language = opts.Language
	} else {
		language = langtype.LangType(filename)
	}

	var input string
	// InputFile > STDIN
	if opts.InputFile != "" {
		inputBS, err := ioutil.ReadFile(opts.InputFile)
		if err != nil {
			println(err.Error())
			return
		}
		input = string(inputBS)
	} else if opts.STDIN {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input += scanner.Text() + "\n"
		}
	}

	client := paizaio.NewClient(paizaio.Config{
		Host:   PAIZA_IO_URL,
		APIKey: PAIZA_IO_API_KEY,
	})

	req1 := &paizaio.RunnersCreateRequest{
		Language:        language,
		SourceCode:      source,
		Input:           input,
		Longpoll:        true,
		LongpollTimeout: 16,
	}

	res1, err := client.RunnersCreate(req1)
	if err != nil {
		println(err.Error())
		return
	}

	req2 := &paizaio.RunnersGetDetailsRequest{
		ID: res1.ID,
	}

	res2, err := client.RunnersGetDetails(req2)
	if err != nil {
		println(err.Error())
		return
	}

	fmt.Fprint(os.Stdout, res2.BuildSTDOUT)
	fmt.Fprint(os.Stderr, res2.BuildSTDERR)

	if res2.BuildExitCode != 0 {
		os.Exit(res2.BuildExitCode)
	}

	fmt.Fprint(os.Stdout, res2.STDOUT)
	fmt.Fprint(os.Stderr, res2.STDERR)

	if opts.Detail {
		fmt.Println()
		if res2.BuildResult != "" {
			fmt.Printf("build_time:   %ss\n", res2.BuildTime)
			fmt.Printf("build_memory: %dKB\n", res2.BuildMemory/1024)
		}
		fmt.Printf("time:         %ss\n", res2.Time)
		fmt.Printf("memory:       %dKB\n", res2.Memory/1024)
	}

	if res2.ExitCode != 0 {
		fmt.Fprintln(os.Stderr, res2.Result)
	}

	os.Exit(int(res2.ExitCode))
}
