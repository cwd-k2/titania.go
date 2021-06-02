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
	Time                string
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

func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *TestCaseResult {
	dirname := filepath.Join(t.Name, target.Name, tcase.Name)
	// TODO: refactoring
	// fire paiza.io API
	sres1 := t.do(dirname, target.Language, target.FileName, []string{tcase.InputFileName})

	var (
		result string
		errstr string = sres1.Error
		// anyting other than stdout
		others []string = []string{sres1.BuildStdoutFileName, sres1.BuildStderrFileName, sres1.StderrFileName}
	)

	// making result string
	if t.TestMethod != nil && sres1.Result == t.TestMethod.OnResult {
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

		if sres2.Result == "SUCCESS" {
			byteArray, _ := ioutil.ReadFile(sres2.StdoutFileName)
			result = strings.TrimRight(string(byteArray), "\n") // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", sres2.Result)
		}

		errstr += sres2.Error
		others = append(others, sres2.BuildStdoutFileName, sres2.BuildStderrFileName, sres2.StderrFileName)

	} else if sres1.Result == "SUCCESS" {
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

	// TODO: how can I treat build time?
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
	if err != nil {
		return handle(err)
	}

	ret := &singleresult{
		StdoutFileName:      stdoutFileName,
		StderrFileName:      stderrFileName,
		BuildStdoutFileName: buildStdoutFileName,
		BuildStderrFileName: buildStderrFileName,
	}

	if res.BuildExitCode != 0 {
		result := fmt.Sprintf("BUILD %s", strings.ToUpper(res.BuildResult))
		ret.Result = result
		ret.Time = res.BuildTime
		return ret
	}

	if res.ExitCode != 0 || res.Result != "success" {
		result := fmt.Sprintf("EXECUTION %s", strings.ToUpper(res.Result))
		ret.Result = result
		ret.Time = res.Time
		return ret
	}

	ret.Result = strings.ToUpper(res.Result)
	ret.Time = res.Time
	return ret
}

func handle(err error) *singleresult {
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

	return &singleresult{Result: result, Error: errstr}
}
