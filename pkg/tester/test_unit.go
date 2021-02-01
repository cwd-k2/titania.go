package tester

import (
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/paizaio"
)

type TestUnit struct {
	Name        string
	Client      *paizaio.Client
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
	view        Viewer
}

// Reads given directory and create an instance of TestUnit.
// if failed to load Config/TestTargets/TestCases, returns nil (no error).
func NewTestUnit(dirname string) *TestUnit {
	basepath, err := filepath.Abs(dirname)
	if err != nil {
		logger.Printf("%+v\n", err)
		return nil
	}

	// 設定
	config := NewConfig(basepath)
	if config == nil {
		return nil
	}

	// paiza.io API クライアント
	client := paizaio.NewClient(config.ClientConfig)

	// ソースコード
	targets := MakeTestTargets(basepath, config.TestTarget)
	if len(targets) == 0 {
		return nil
	}

	// テストケース
	tcases := MakeTestCases(basepath, config.TestCase)
	if len(tcases) == 0 {
		return nil
	}

	// テストメソッド
	tmethod := NewTestMethod(basepath, config.TestMethod)

	// Viewer
	indices := make([]string, 0)
	for _, target := range targets {
		indices = append(indices, target.Name)
	}
	view := NewView(dirname, len(targets), len(tcases), indices)

	return &TestUnit{dirname, client, tmethod, targets, tcases, view}
}

// Execute test (itself) using paiza.io API.
// Any errors are included in returning values.
func (t *TestUnit) Exec() *Outcome {
	curr := 0
	stop := len(t.TestTargets) * len(t.TestCases)

	fruits := make([]*Fruit, len(t.TestTargets))
	for i, target := range t.TestTargets {
		fruits[i] = &Fruit{target.Name, target.Language, target.Expect, make([]*Detail, len(t.TestCases))}
	}

	// idiom: sending multiple value with a single channel
	ch := make(chan func() (int, int, *Detail), stop)
	fn := func(i, j int, target *TestTarget, tcase *TestCase) {
		detail := t.exec(target, tcase)
		ch <- func() (int, int, *Detail) { return i, j, detail }
	}

	t.view.Init()

	// Each test is executed asynchronously
	for i, target := range t.TestTargets {
		for j, tcase := range t.TestCases {
			go fn(i, j, target, tcase)
		}
	}

	for res := range ch {
		i, j, d := res()

		fruits[i].Details[j] = d

		t.view.Update(i)

		if curr++; curr == stop {
			close(ch)
		}
	}

	t.view.Done()

	outcome := &Outcome{t.Name, "default", fruits}
	if t.TestMethod != nil {
		outcome.TestMethod = t.TestMethod.Name
	}

	return outcome
}
