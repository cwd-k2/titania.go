package tester

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/file"
	"github.com/cwd-k2/titania.go/pkg/runner"
)

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
	others := []string{}

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
	others = append(
		others,
		apires.BuildStdoutFile,
		apires.BuildStderrFile,
		apires.StderrFile,
	)

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
				apires.WriteStdoutTo(inputs)
			case "stderr":
				apires.WriteStderrTo(inputs)
			case "build_stdout":
				apires.WriteBuildStdoutTo(inputs)
			case "build_stderr":
				apires.WriteBuildStderrTo(inputs)
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
			result, _ = res.ReadTestResult() // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", res.Result)
		}

		// append information to the result
		errstr += res.Error
		others = append(
			others,
			res.BuildStdoutFile,
			res.BuildStderrFile,
			res.StdoutFile,
			res.StderrFile,
		)

	} else if apires.BuildExitCode == 0 && apires.ExitCode == 0 {
		// simple comparison of answer and output (byte level)
		if b, _ := file.Equal(apires.StdoutFile, tc.AnswerFile); b {
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
		Output: apires.StdoutFile,
		Others: others,
		Errors: errstr,
	}
}

// fire paiza.io api and convert the response to a shape
func (t *TestUnit) do(name, language string, source, input io.Reader) *apiresult {
	// the real entities where to write api result outputs
	dirname := filepath.Join(tmpdir, name)
	os.MkdirAll(dirname, 0755)

	stdoutFile := filepath.Join(dirname, "stdout")
	stderrFile := filepath.Join(dirname, "stderr")
	buildStdoutFile := filepath.Join(dirname, "build_stdout")
	buildStderrFile := filepath.Join(dirname, "build_stderr")

	stdoutFp, _ := os.Create(stdoutFile)
	stderrFp, _ := os.Create(stderrFile)
	buildStdoutFp, _ := os.Create(buildStdoutFile)
	buildStderrFp, _ := os.Create(buildStderrFile)

	defer stdoutFp.Close()
	defer stderrFp.Close()
	defer buildStdoutFp.Close()
	defer buildStderrFp.Close()

	stdoutEnt := bufio.NewWriter(stdoutFp)
	stderrEnt := bufio.NewWriter(stderrFp)
	buildStdoutEnt := bufio.NewWriter(buildStdoutFp)
	buildStderrEnt := bufio.NewWriter(buildStderrFp)

	defer stdoutEnt.Flush()
	defer stderrEnt.Flush()
	defer buildStdoutEnt.Flush()
	defer buildStderrEnt.Flush()

	// fire paiza.io API
	res, err := t.Runner.Run(&runner.OrderSpec{
		Language:    language,
		SourceCode:  source,
		Input:       input,
		Stdout:      stdoutEnt,
		Stderr:      stderrEnt,
		BuildStdout: buildStdoutEnt,
		BuildStderr: buildStderrEnt,
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
		StdoutFile:      stdoutFile,
		StderrFile:      stderrFile,
		BuildStdoutFile: buildStdoutFile,
		BuildStderrFile: buildStderrFile,
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
