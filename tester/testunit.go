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
	SourceCodes []*SourceCode
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
	client := new(client.Client)
	client.Host = config.Host
	client.APIKey = config.APIKey

	// ソースコード
	sourceCodes := MakeSourceCode(basepath, languages, config.SourceCodeDirectories)

	// ソースコードがなければ実行しない
	if len(sourceCodes) == 0 {
		return nil
	}

	// テストケース
	testCases := MakeTestCases(
		basepath,
		config.TestCaseDirectories,
		config.TestCaseInputExtension,
		config.TestCaseAnswerExtension)

	// テストケースがなければ実行しない
	if len(testCases) == 0 {
		return nil
	}

	// テストメソッド
	testMethod := NewTestMethod(basepath, config.TestMethodFileName)

	testUnit := new(TestUnit)
	testUnit.Name = dirname
	testUnit.Client = client
	testUnit.Config = config
	testUnit.SourceCodes = sourceCodes
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
	fruits := make([]*Fruit, len(testUnit.SourceCodes))

	for i, sourceCode := range testUnit.SourceCodes {
		fruit := new(Fruit)
		fruit.SourceCode = sourceCode.Name
		fruit.Language = sourceCode.Language
		fruit.Details = make([]*Detail, len(testUnit.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	testUnit.goEach(func(i, j int, sourceCode *SourceCode, testCase *TestCase) {
		detail := testUnit.exec(sourceCode, testCase)
		fruits[i].Details[j] = detail
		ch <- i
	})

	curr := 0
	stop := len(testUnit.SourceCodes) * len(testUnit.TestCases)
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

func (testUnit *TestUnit) exec(sourceCode *SourceCode, testCase *TestCase) *Detail {

	detail := sourceCode.Exec(testUnit.Client, testCase)

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

	return detail

}

func (testUnit *TestUnit) goEach(delegateFunc func(int, int, *SourceCode, *TestCase)) {
	for i, sourceCode := range testUnit.SourceCodes {
		for j, testCase := range testUnit.TestCases {
			go delegateFunc(i, j, sourceCode, testCase)
		}
	}
}
