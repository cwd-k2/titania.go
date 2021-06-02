package tester

type TestUnitResult struct {
	Name        string
	TestMethod  string
	TestTargets []*TestTargetResult
}

type TestTargetResult struct {
	Name      string
	Language  string
	Expect    string
	TestCases []*TestCaseResult
}

type TestCaseResult struct {
	Name       string
	Result     string
	IsExpected bool
	Time       string
	Output     string
	Others     []string
	Error      string
}
