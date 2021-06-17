package tester

import (
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
	Result          string
	BuildTime       string
	BuildExitCode   int
	Time            string
	ExitCode        int
	BuildStdoutData []byte
	BuildStderrData []byte
	StdoutData      []byte
	StderrData      []byte
	Error           string
}

// the first token to be the result
func readResult(data []byte) string {
	str := string(data)
	if str == "" {
		return "<NONE>"
	} else {
		return strings.TrimRight(strings.SplitN(str, "\n", 1)[0], "\n")
	}
}

func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *TestCaseResult {
	var (
		result string
		errstr string
	)
	// TODO: refactoring
	dirname := filepath.Join(t.Name, target.Name, tcase.Name)

	source := bytes.NewBuffer([]byte{})
	inputs := bytes.NewBuffer([]byte{})
	others := bytes.NewBuffer([]byte{})

	// fire paiza.io API
	source.Write(target.CodeData)
	inputs.Write(tcase.InputData)
	sres1 := t.do(dirname, target.Language, source, inputs)
	// anyting other than stdout
	errstr += sres1.Error
	others.Write(sres1.BuildStdoutData)
	others.Write(sres1.BuildStderrData)
	others.Write(sres1.StderrData)

	expect, ok := target.Expect[tcase.Name]
	if !ok {
		expect = target.Expect["default"]
	}

	// making result string
	if t.TestMethod != nil && (sres1.ExitCode == t.TestMethod.OnExit || sres1.BuildExitCode == t.TestMethod.OnExit-512) {
		source := bytes.NewBuffer([]byte{})
		source.Write(t.TestMethod.CodeData)

		inputs := bytes.NewBuffer([]byte{})
		for _, what := range t.TestMethod.InputOrder {
			switch what {
			case "input":
				inputs.Write(tcase.InputData)
			case "answer":
				inputs.Write(tcase.AnswerData)
			case "source_code":
				inputs.Write(target.CodeData)
			case "stdout":
				inputs.Write(sres1.StdoutData)
			case "stderr":
				inputs.Write(sres1.StderrData)
			case "build_stdout":
				inputs.Write(sres1.BuildStdoutData)
			case "build_stderr":
				inputs.Write(sres1.BuildStderrData)
			case "delimiter":
				inputs.WriteString(t.TestMethod.Delimiter)
			case "newline":
				inputs.WriteString("\n")
			case "tab":
				inputs.WriteString("\t")
			}
		}
		// TestMethod
		sres2 := t.do(filepath.Join(dirname, t.TestMethod.Name), t.TestMethod.Language, source, inputs)

		// TestMethod should gracefully terminate.
		if sres2.BuildExitCode == 0 && sres2.ExitCode == 0 {
			result = readResult(sres2.StdoutData) // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", sres2.Result)
		}

		errstr += sres2.Error
		others.Write(sres2.BuildStdoutData)
		others.Write(sres2.BuildStderrData)
		others.Write(sres2.StdoutData)
		others.Write(sres2.StderrData)

	} else if sres1.BuildExitCode == 0 && sres1.ExitCode == 0 {
		// simple comparison
		if bytes.Equal(sres1.StdoutData, tcase.AnswerData) {
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
		IsExpected: result == expect,
		Time:       sres1.Time,
		Output:     string(sres1.StdoutData),
		Others:     others.String(),
		Error:      errstr,
	}
}

func (t *TestUnit) do(name, language string, source, input io.Reader) *singleresult {
	// power is power
	stdoutEnt := ioutil.Discard
	stderrEnt := ioutil.Discard
	buildStdoutEnt := ioutil.Discard
	buildStderrEnt := ioutil.Discard

	stdoutBuf := bytes.NewBuffer([]byte{})
	stderrBuf := bytes.NewBuffer([]byte{})
	buildStdoutBuf := bytes.NewBuffer([]byte{})
	buildStderrBuf := bytes.NewBuffer([]byte{})

	if tmpdir != "" {
		dirname := filepath.Join(tmpdir, name)
		os.MkdirAll(dirname, 0755)

		stdoutFp, _ := os.Create(filepath.Join(dirname, "stdout"))
		stderrFp, _ := os.Create(filepath.Join(dirname, "stderr"))
		buildStdoutFp, _ := os.Create(filepath.Join(dirname, "build_stdout"))
		buildStderrFp, _ := os.Create(filepath.Join(dirname, "build_stderr"))

		stdoutEnt = stdoutFp
		stderrEnt = stderrFp
		buildStdoutEnt = buildStdoutFp
		buildStderrEnt = buildStderrFp

		defer stdoutFp.Close()
		defer stderrFp.Close()
		defer buildStdoutFp.Close()
		defer buildStderrFp.Close()
	}

	res, err := t.Runner.Run(&runner.OrderSpec{
		Language:    language,
		SourceCode:  source,
		Input:       input,
		Stdout:      io.MultiWriter(stdoutEnt, stdoutBuf),
		Stderr:      io.MultiWriter(stderrEnt, stderrBuf),
		BuildStdout: io.MultiWriter(buildStdoutEnt, buildStdoutBuf),
		BuildStderr: io.MultiWriter(buildStderrEnt, buildStderrBuf),
	})

	ret := &singleresult{
		Time:            res.Time,
		ExitCode:        res.ExitCode,
		BuildTime:       res.BuildTime,
		BuildExitCode:   res.BuildExitCode,
		StdoutData:      stdoutBuf.Bytes(),
		StderrData:      stderrBuf.Bytes(),
		BuildStdoutData: buildStdoutBuf.Bytes(),
		BuildStderrData: buildStderrBuf.Bytes(),
	}

	if err != nil {
		ret.Result, ret.Error = handle(err)
		ret.BuildExitCode, ret.ExitCode = -1, -1
		return ret
	}

	// TIMEOUT exit code to be 124 (like coreutils timeout)
	if res.BuildResult == "timeout" {
		ret.BuildExitCode = 124
	}
	if res.Result == "timeout" {
		ret.ExitCode = 124
	}

	if ret.BuildExitCode != 0 {
		ret.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(res.BuildResult))
	} else if ret.ExitCode != 0 {
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
