package tester

import (
	"log"
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type TestUnit struct {
	Name        string
	Runner      *runner.Runner
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
	view        Viewer
}

// Reads given directory and create an instance of TestUnit.
// if failed to load Config/TestTargets/TestCases, returns nil (no error).
func ReadTestUnit(dirname string, config *Config) *TestUnit {
	basepath, err := filepath.Abs(dirname)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil
	}

	// Code Runner
	runner := runner.NewRunner(config.ClientConfig)

	// ソースコード
	targets := ReadTestTargets(basepath, config.TestTarget)
	if len(targets) == 0 {
		return nil
	}

	// テストケース
	tcases := ReadTestCases(basepath, config.TestCase)
	if len(tcases) == 0 {
		return nil
	}

	// テストメソッド
	tmethod := ReadTestMethod(basepath, config.TestMethod)

	// Viewer
	indices := make([]string, 0)
	for _, target := range targets {
		indices = append(indices, target.Name)
	}
	view := NewView(dirname, len(targets), len(tcases), indices)

	return &TestUnit{dirname, runner, tmethod, targets, tcases, view}
}

// Execute test (itself) using paiza.io API.
// Any errors are included in returning values.
func (t *TestUnit) Exec() *TestUnitResult {
	curr := 0
	stop := len(t.TestTargets) * len(t.TestCases)

	tresults := make([]*TestTargetResult, len(t.TestTargets))
	for i, target := range t.TestTargets {
		tresults[i] = &TestTargetResult{target.Name, target.Language, target.Expect, make([]*TestCaseResult, len(t.TestCases))}
	}

	// idiom: sending multiple value with a single channel
	ch := make(chan func() (int, int, *TestCaseResult), stop)
	fn := func(i, j int, target *TestTarget, tcase *TestCase) {
		cresult := t.exec(target, tcase)
		ch <- func() (int, int, *TestCaseResult) { return i, j, cresult }
	}

	t.view.Init()

	// Each test is executed asynchronously
	for i, target := range t.TestTargets {
		for j, tcase := range t.TestCases {
			go fn(i, j, target, tcase)
		}
	}

	for res := range ch {
		i, j, cresult := res()

		tresults[i].TestCases[j] = cresult

		t.view.Update(i)

		if curr++; curr == stop {
			close(ch)
		}
	}

	t.view.Done()

	uresult := &TestUnitResult{t.Name, "default", tresults}
	if t.TestMethod != nil {
		uresult.TestMethod = t.TestMethod.Name
	}

	return uresult
}
