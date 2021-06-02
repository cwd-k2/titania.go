package tester

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type singleresult struct {
	Result              string
	BuildTime           string
	BuildExitCode       int
	Time                string
	ExitCode            int
	BuildStdoutFileName string
	BuildStderrFileName string
	StdoutFileName      string
	StderrFileName      string
	Error               string
}

func sameFileContents(fname1, fname2 string) bool {
	// mmm...
	b1, _ := ioutil.ReadFile(fname1)
	b2, _ := ioutil.ReadFile(fname2)

	return bytes.Equal(b1, b2)
}

// the first token to be the result
func readResult(filename string) string {
	fp, _ := os.Open(filename)
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	if scanner.Scan() {
		return scanner.Text()
	} else {
		return "<NONE>"
	}
}

func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *TestCaseResult {
	var (
		result string
		errstr string
	)
	// TODO: refactoring
	dirname := filepath.Join(t.Name, target.Name, tcase.Name)
	// fire paiza.io API
	sres1 := t.do(dirname, target.Language, target.FileName, []string{tcase.InputFileName})
	// anyting other than stdout
	others := []string{sres1.BuildStdoutFileName, sres1.BuildStderrFileName, sres1.StderrFileName}

	// making result string
	if t.TestMethod != nil && (sres1.ExitCode == t.TestMethod.OnExit || sres1.BuildExitCode == t.TestMethod.OnExit-512) {
		// ex: stdin + '\000' + stdout + '\000' + answer
		var inputs []string
		for _, what := range t.TestMethod.InputOrder {
			switch what {
			case "input":
				inputs = append(inputs, tcase.InputFileName)
			case "answer":
				inputs = append(inputs, tcase.AnswerFileName)
			case "source_code":
				inputs = append(inputs, target.FileName)
			case "stdout":
				inputs = append(inputs, sres1.StdoutFileName)
			case "stderr":
				inputs = append(inputs, sres1.StderrFileName)
			case "build_stdout":
				inputs = append(inputs, sres1.BuildStdoutFileName)
			case "build_stderr":
				inputs = append(inputs, sres1.BuildStderrFileName)
			}
		}
		// TestMethod
		sres2 := t.do(filepath.Join(dirname, t.TestMethod.Name), t.TestMethod.Language, t.TestMethod.FileName, inputs)

		// TestMethod should gracefully terminate.
		if sres2.BuildExitCode == 0 && sres2.ExitCode == 0 {
			result = readResult(sres2.StdoutFileName) // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", sres2.Result)
		}

		errstr += sres2.Error
		others = append(others, sres2.BuildStdoutFileName, sres2.BuildStderrFileName, sres1.StdoutFileName, sres2.StderrFileName)

	} else if sres1.BuildExitCode == 0 && sres1.ExitCode == 0 {
		// simple comparison
		if sameFileContents(sres1.StdoutFileName, tcase.AnswerFileName) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	} else {
		result = sres1.Result
	}

	return &TestCaseResult{
		Name:       tcase.Name,
		Result:     result,
		IsExpected: result == target.Expect,
		Time:       sres1.Time,
		Output:     sres1.StdoutFileName,
		Others:     others,
		Error:      errstr,
	}
}

func (t *TestUnit) do(name, language, sourceFileName string, inputFileNames []string) *singleresult {
	// power is power
	dirname := filepath.Join(tmpdir, name)
	os.MkdirAll(dirname, 0755)

	sourcecode, _ := os.Open(sourceFileName)
	defer sourcecode.Close()

	inputs := make([]io.Reader, 0)
	for _, inputFileName := range inputFileNames {
		input, _ := os.Open(inputFileName)
		inputs = append(inputs, bufio.NewReader(input))
		defer input.Close()
	}

	stdoutFileName := filepath.Join(dirname, "stdout")
	stdout, _ := os.OpenFile(stdoutFileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer stdout.Close()
	stderrFileName := filepath.Join(dirname, "stderr")
	stderr, _ := os.OpenFile(stderrFileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer stderr.Close()
	buildStdoutFileName := filepath.Join(dirname, "build_stdout")
	buildstdout, _ := os.OpenFile(buildStdoutFileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer buildstdout.Close()
	buildStderrFileName := filepath.Join(dirname, "build_stderr")
	buildstderr, _ := os.OpenFile(buildStderrFileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer buildstderr.Close()

	res, err := t.Runner.Run(&runner.OrderSpec{
		Language:       language,
		SourceCode:     bufio.NewReader(sourcecode),
		Inputs:         inputs,
		InputDelimiter: "\x00",
		Stdout:         stdout,
		Stderr:         stderr,
		BuildStdout:    buildstdout,
		BuildStderr:    buildstderr,
	})

	ret := &singleresult{
		Time:                res.Time,
		ExitCode:            res.ExitCode,
		BuildTime:           res.BuildTime,
		BuildExitCode:       res.BuildExitCode,
		StdoutFileName:      stdoutFileName,
		StderrFileName:      stderrFileName,
		BuildStdoutFileName: buildStdoutFileName,
		BuildStderrFileName: buildStderrFileName,
	}

	if err != nil {
		ret.Result, ret.Error = handle(err)
		ret.BuildExitCode, ret.ExitCode = -1, -1
		return ret
	}

	if res.BuildExitCode != 0 {
		ret.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(res.BuildResult))
	} else if res.ExitCode != 0 {
		ret.Result = fmt.Sprintf("EXECUTION %s", strings.ToUpper(res.Result))
	} else {
		ret.Result = strings.ToUpper(res.Result)
	}
	return ret
}

func handle(err error) (string, string) {
	var result, errstr string

	switch err := err.(type) {
	case runner.ServerError:
		result = "SERVER ERROR"
		errstr = fmt.Sprintf("HTTP response status code: %d\n%s", err.Code, err.Error())
	case runner.ClientError:
		result = "CLIENT ERROR"
		errstr = fmt.Sprintf("HTTP response status code: %d\n%s", err.Code, err.Error())
	case runner.RunnerError:
		result = "RUNNER ERROR"
		errstr = fmt.Sprintf("Error occurred at paiza.io code runner.\n%s", err.Error())
	default:
		result = "TESTER ERROR"
		errstr = err.Error()
	}

	return result, errstr
}
