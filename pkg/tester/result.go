package tester

import (
	"bufio"
	"io"
	"os"
)

type TestUnitResult struct {
	Name        string
	TestMethod  string
	TestTargets []*TestTargetResult
}

type TestTargetResult struct {
	Name      string
	Language  string
	TestCases []*TestCaseResult
}

type TestCaseResult struct {
	Name   string
	Time   string
	Expect string
	Result string
	Output string
	Others []string
	Errors string
}

// different from the structs above; used internally
type apiresult struct {
	Result          string
	BuildTime       string
	BuildExitCode   int
	Time            string
	ExitCode        int
	BuildStdoutFile string
	BuildStderrFile string
	StdoutFile      string
	StderrFile      string
	Error           string
}

// check if the result's exit code is expected by test method
func isExpectedExitCode(r *apiresult, t *TestMethod) bool {
	return r.ExitCode == t.OnExit || r.BuildExitCode == t.OnExit-512
}

// read the first line
func (r *apiresult) ReadTestResult() (string, error) {
	fp, err := os.Open(r.StdoutFile)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	tk, err := bufio.NewReader(fp).ReadString('\n')
	if err != nil && err == io.EOF {
		return "<NONE>", nil
	} else if err != nil {
		return "", err
	} else {
		return tk, nil
	}
}

// write stdout to the writer
func (r *apiresult) WriteStdoutTo(w io.Writer) error {
	fp, err := os.Open(r.StdoutFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}

// write stderr to the writer
func (r *apiresult) WriteStderrTo(w io.Writer) error {
	fp, err := os.Open(r.StderrFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}

// write build_stdout to the writer
func (r *apiresult) WriteBuildStdoutTo(w io.Writer) error {
	fp, err := os.Open(r.BuildStdoutFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}

// write build_stderr to the writer
func (r *apiresult) WriteBuildStderrTo(w io.Writer) error {
	fp, err := os.Open(r.BuildStderrFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}
