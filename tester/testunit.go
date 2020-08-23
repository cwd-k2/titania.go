package tester

import (
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// TestUnit
// contains paiza.io API client, config, SourceCodes, and TestCases
// physically, this stands for a directory.
type TestUnit struct {
	Name        string
	Client      *client.Client
	Config      *Config
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
}

// NewTestUnit
// returns *TestUnit
func NewTestUnit(dirname string, languages []string) *TestUnit {
	basepath, err := filepath.Abs(dirname)
	// ここのエラーは公式のドキュメント見てもわからんのだけど何？
	if err != nil {
		println(err)
		return nil
	}

	// 設定
	config := NewConfig(basepath)
	if config == nil {
		return nil
	}

	// paiza.io API クライアント
	client := client.NewClient(config.ClientConfig)

	// ソースコード
	testTargets := MakeTestTargets(basepath, languages, config.TestTarget)

	// ソースコードがなければ実行しない
	if len(testTargets) == 0 {
		return nil
	}

	// テストケース
	testCases := MakeTestCases(basepath, config.TestCase)

	// テストケースがなければ実行しない
	if len(testCases) == 0 {
		return nil
	}

	// テストメソッド
	testMethod := NewTestMethod(basepath, config.TestMethod)

	testUnit := new(TestUnit)
	testUnit.Name = dirname
	testUnit.Client = client
	testUnit.Config = config
	testUnit.TestTargets = testTargets
	testUnit.TestCases = testCases
	testUnit.TestMethod = testMethod

	return testUnit
}

func MakeTestUnits(directories, languages []string) []*TestUnit {
	testUnits := make([]*TestUnit, 0, len(directories))
	for _, dirname := range directories {
		testUnit := NewTestUnit(dirname, languages)
		if testUnit != nil {
			testUnits = append(testUnits, testUnit)
		}
	}
	return testUnits
}

func (testUnit *TestUnit) Exec(view View) *Outcome {
	ch := make(chan int)
	fruits := make([]*Fruit, len(testUnit.TestTargets))

	for i, testTarget := range testUnit.TestTargets {
		fruit := new(Fruit)
		fruit.TestTarget = testTarget.Name
		fruit.Language = testTarget.Language
		fruit.Expect = testTarget.Expect
		fruit.Details = make([]*Detail, len(testUnit.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	testUnit.goEach(func(i, j int, testTarget *TestTarget, testCase *TestCase) {
		detail := testUnit.exec(testTarget, testCase)
		fruits[i].Details[j] = detail
		ch <- i
	})

	curr := 0
	stop := len(testUnit.TestTargets) * len(testUnit.TestCases)
	for i := range ch {
		curr++
		view.Update(i)
		if curr == stop {
			close(ch)
		}
	}

	outcome := new(Outcome)
	outcome.TestUnit = testUnit.Name
	if testUnit.TestMethod != nil {
		outcome.TestMethod = testUnit.TestMethod.Name
	} else {
		outcome.TestMethod = "default"
	}
	outcome.Fruits = fruits

	return outcome
}

func (testUnit *TestUnit) exec(testTarget *TestTarget, testCase *TestCase) *Detail {

	detail := testTarget.Exec(testUnit.Client, testCase)

	// if result not set (this means execution was successful).
	if detail.Result == "" {

		if testUnit.TestMethod != nil {
			// use custom testing method.
			res, ers := testUnit.TestMethod.Exec(testUnit.Client, testCase, detail)
			detail.Result = strings.TrimRight(res, "\n")
			detail.Error += ers
		} else {
			// just compare ouput and expected answer.
			if detail.Output == testCase.Answer {
				detail.Result = "PASS"
			} else {
				detail.Result = "FAIL"
			}
		}

	}

	detail.Expected = detail.Result == testTarget.Expect

	return detail

}

func (testUnit *TestUnit) goEach(delegateFunc func(int, int, *TestTarget, *TestCase)) {
	for i, testTarget := range testUnit.TestTargets {
		for j, testCase := range testUnit.TestCases {
			go delegateFunc(i, j, testTarget, testCase)
		}
	}
}
