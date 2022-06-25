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

type apiresult struct {
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

// check if the result's exit code is expected by test method
func isExpectedExitCode(r *apiresult, t *TestMethod) bool {
	return r.ExitCode == t.OnExit || r.BuildExitCode == t.OnExit-512
}

// the first token to be the result
func readTestResult(data []byte) string {
	if len(data) == 0 {
		return "<NONE>"
	}
	return string(bytes.SplitN(data, []byte{'\n'}, 2)[0])
}

// Execute the 'test' using paiza.IO api
func (t *TestUnit) exec(i, j int) *TestCaseResult {
	// TODO: refactoring
	var (
		result string
		errstr string
	)

	tt := t.TestTargets[i]
	tc := t.TestCases[j]
	tm := t.TestMethod

	dirname := filepath.Join(t.Name, tt.Name, tc.Name)

	source := bytes.NewBuffer([]byte{})
	inputs := bytes.NewBuffer([]byte{})
	others := bytes.NewBuffer([]byte{})

	expect, ok := tt.Expect[tc.Name]
	if !ok {
		expect = tt.Expect["default"]
	}

	tt.WriteSouceCodeTo(source)
	tc.WriteInputDataTo(inputs)
	// fire paiza.io API
	apires := t.do(dirname, tt.Language, source, inputs)
	// append information other than stdout
	errstr += apires.Error
	others.Write(apires.BuildStdoutData)
	others.Write(apires.BuildStderrData)
	others.Write(apires.StderrData)

	// making result string
	if tm != nil && isExpectedExitCode(apires, tm) {
		// if test method exists, fire paiza.io api again to test the result
		source := bytes.NewBuffer([]byte{})
		inputs := bytes.NewBuffer([]byte{})
		// write test method code
		tm.WriteSouceCodeTo(source)
		// create input for test method
		for _, what := range tm.InputOrder {
			switch what {
			case "input":
				tc.WriteInputDataTo(inputs)
			case "answer":
				tc.WriteAnswerDataTo(inputs)
			case "source_code":
				tt.WriteSouceCodeTo(inputs)
			case "stdout":
				inputs.Write(apires.StdoutData)
			case "stderr":
				inputs.Write(apires.StderrData)
			case "build_stdout":
				inputs.Write(apires.BuildStdoutData)
			case "build_stderr":
				inputs.Write(apires.BuildStderrData)
			case "language":
				inputs.WriteString(tt.Language)
			case "delimiter":
				inputs.WriteString(tm.Delimiter)
			case "newline":
				inputs.WriteString("\n")
			case "tab":
				inputs.WriteString("\t")
			}
		}
		// execute TestMethod
		res := t.do(filepath.Join(dirname, tm.Name), tm.Language, source, inputs)

		// TestMethod should gracefully terminate.
		if res.BuildExitCode == 0 && res.ExitCode == 0 {
			result = readTestResult(res.StdoutData) // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", res.Result)
		}

		// append information to the result
		errstr += res.Error
		others.Write(res.BuildStdoutData)
		others.Write(res.BuildStderrData)
		others.Write(res.StdoutData)
		others.Write(res.StderrData)

	} else if apires.BuildExitCode == 0 && apires.ExitCode == 0 {
		// read out the expected answer
		answerBuf := bytes.NewBuffer([]byte{})
		tc.WriteAnswerDataTo(answerBuf)
		// simple comparison of answer and output (byte level)
		if bytes.Equal(apires.StdoutData, answerBuf.Bytes()) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	} else {
		result = apires.Result
	}

	return &TestCaseResult{
		Name:   tc.Name,
		Time:   apires.Time,
		Expect: expect,
		Result: result,
		Output: bytes.NewReader(apires.StdoutData),
		Others: others,
		Errors: errstr,
	}
}

// fire paiza.io api and convert the response to a shape
func (t *TestUnit) do(name, language string, source, input io.Reader) *apiresult {
	// the real entities where to write api result outputs
	// default: discard
	stdoutEnt := io.Discard
	stderrEnt := io.Discard
	buildStdoutEnt := io.Discard
	buildStderrEnt := io.Discard
	// if output directory is specified, then create files and set them
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
	// bytes.NewBuffer takes the ownership of the byte slice
	stdoutBuf := bytes.NewBuffer([]byte{})
	stderrBuf := bytes.NewBuffer([]byte{})
	buildStdoutBuf := bytes.NewBuffer([]byte{})
	buildStderrBuf := bytes.NewBuffer([]byte{})

	// fire paiza.io API
	res, err := t.Runner.Run(&runner.OrderSpec{
		Language:    language,
		SourceCode:  source,
		Input:       input,
		Stdout:      io.MultiWriter(stdoutEnt, stdoutBuf),
		Stderr:      io.MultiWriter(stderrEnt, stderrBuf),
		BuildStdout: io.MultiWriter(buildStdoutEnt, buildStdoutBuf),
		BuildStderr: io.MultiWriter(buildStderrEnt, buildStderrBuf),
	})

	// something went wrong?
	if err != nil {
		result, errstr := handle(err)
		return &apiresult{
			Time:          "-1",
			Result:        result,
			Error:         errstr,
			ExitCode:      -1,
			BuildExitCode: -1,
		}
	}

	// creating returned object: apiResult
	ret := &apiresult{
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
	// create result string according to the exit codes
	if ret.BuildExitCode != 0 {
		ret.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(res.BuildResult))
	} else if ret.ExitCode != 0 {
		ret.Result = fmt.Sprintf("EXECUTION %s", strings.ToUpper(res.Result))
	} else {
		ret.Result = strings.ToUpper(res.Result)
	}

	return ret
}

// create result and excuse strings from an error
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
