package tester

import "io"

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
	Output io.Reader
	Others io.Reader
	Errors string
}
