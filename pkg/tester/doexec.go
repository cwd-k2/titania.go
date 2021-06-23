package tester

import (
	"bytes"
	"fmt"
	"io"
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
		return strings.Split(str, "\n")[0]
	}
}

func (t *TestUnit) exec(i, j int) *TestCaseResult {
	// TODO: refactoring
	var (
		result string
		errstr string
	)

	dirname := filepath.Join(t.Name, t.TestTargets[i].Name, t.TestCases[j].Name)

	source := bytes.NewBuffer([]byte{})
	inputs := bytes.NewBuffer([]byte{})
	others := bytes.NewBuffer([]byte{})

	expect, ok := t.TestTargets[i].Expect[t.TestCases[j].Name]
	if !ok {
		expect = t.TestTargets[i].Expect["default"]
	}

	source.Write(t.TestTargets[i].CodeData)
	inputs.Write(t.TestCases[j].InputData)
	// fire paiza.io API
	ttsres := t.do(dirname, t.TestTargets[i].Language, source, inputs)
	errstr += ttsres.Error
	// anyting other than stdout
	others.Write(ttsres.BuildStdoutData)
	others.Write(ttsres.BuildStderrData)
	others.Write(ttsres.StderrData)

	// making result string
	if t.TestMethod != nil && (ttsres.ExitCode == t.TestMethod.OnExit || ttsres.BuildExitCode == t.TestMethod.OnExit-512) {
		source := bytes.NewBuffer([]byte{})
		source.Write(t.TestMethod.CodeData)
		inputs := bytes.NewBuffer([]byte{})
		for _, what := range t.TestMethod.InputOrder {
			switch what {
			case "input":
				inputs.Write(t.TestCases[j].InputData)
			case "answer":
				inputs.Write(t.TestCases[j].AnswerData)
			case "source_code":
				inputs.Write(t.TestTargets[i].CodeData)
			case "stdout":
				inputs.Write(ttsres.StdoutData)
			case "stderr":
				inputs.Write(ttsres.StderrData)
			case "build_stdout":
				inputs.Write(ttsres.BuildStdoutData)
			case "build_stderr":
				inputs.Write(ttsres.BuildStderrData)
			case "language":
				inputs.WriteString(t.TestTargets[i].Language)
			case "delimiter":
				inputs.WriteString(t.TestMethod.Delimiter)
			case "newline":
				inputs.WriteString("\n")
			case "tab":
				inputs.WriteString("\t")
			}
		}
		// TestMethod
		tmsres := t.do(filepath.Join(dirname, t.TestMethod.Name), t.TestMethod.Language, source, inputs)

		// TestMethod should gracefully terminate.
		if tmsres.BuildExitCode == 0 && tmsres.ExitCode == 0 {
			result = readResult(tmsres.StdoutData) // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", tmsres.Result)
		}

		errstr += tmsres.Error
		others.Write(tmsres.BuildStdoutData)
		others.Write(tmsres.BuildStderrData)
		others.Write(tmsres.StdoutData)
		others.Write(tmsres.StderrData)

	} else if ttsres.BuildExitCode == 0 && ttsres.ExitCode == 0 {
		// simple comparison
		if bytes.Equal(ttsres.StdoutData, t.TestCases[j].AnswerData) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	} else {
		result = ttsres.Result
	}

	return &TestCaseResult{
		Name:   t.TestCases[j].Name,
		Time:   ttsres.Time,
		Expect: expect,
		Result: result,
		Output: bytes.NewReader(ttsres.StdoutData),
		Others: others,
		Errors: errstr,
	}
}

func (t *TestUnit) do(name, language string, source, input io.Reader) *singleresult {
	// power is power
	stdoutEnt := io.Discard
	stderrEnt := io.Discard
	buildStdoutEnt := io.Discard
	buildStderrEnt := io.Discard

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

	if err != nil {
		result, errstr := handle(err)
		return &singleresult{
			Time:          "-1",
			Result:        result,
			Error:         errstr,
			ExitCode:      -1,
			BuildExitCode: -1,
		}
	}

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
