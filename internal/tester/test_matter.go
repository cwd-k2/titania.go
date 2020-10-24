package tester

import (
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/internal/client"
)

// TestMatter
// contains paiza.io API client, SourceCodes, and TestCases
// physically, this stands for a directory.
type TestMatter struct {
	Name        string
	Client      *client.Client
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
}

// NewtestMatter
// returns *testMatter
func NewTestMatter(dirname string, languages []string) *TestMatter {
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

	testMatter := new(TestMatter)
	testMatter.Name = dirname
	testMatter.Client = client
	testMatter.TestTargets = testTargets
	testMatter.TestCases = testCases
	testMatter.TestMethod = testMethod

	return testMatter
}

func MakeTestMatters(directories, languages []string) []*TestMatter {
	testMatters := make([]*TestMatter, 0, len(directories))
	for _, dirname := range directories {
		testMatter := NewTestMatter(dirname, languages)
		if testMatter != nil {
			testMatters = append(testMatters, testMatter)
		}
	}
	return testMatters
}

func (testMatter *TestMatter) Exec(view View) *Outcome {
	curr := 0
	stop := len(testMatter.TestTargets) * len(testMatter.TestCases)

	ch := make(chan int, stop)
	fruits := make([]*Fruit, len(testMatter.TestTargets))

	outcome := new(Outcome)
	outcome.TestMatter = testMatter.Name
	if testMatter.TestMethod != nil {
		outcome.TestMethod = testMatter.TestMethod.Name
	} else {
		outcome.TestMethod = "default"
	}

	for i, testTarget := range testMatter.TestTargets {
		fruit := new(Fruit)
		fruit.TestTarget = testTarget.Name
		fruit.Language = testTarget.Language
		fruit.Expect = testTarget.Expect
		fruit.Details = make([]*Detail, len(testMatter.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	testMatter.goEach(func(i, j int, testTarget *TestTarget, testCase *TestCase) {
		fruits[i].Details[j] = testMatter.exec(testTarget, testCase)
		ch <- i
	})

	for i := range ch {
		curr++
		view.Update(i)
		if curr == stop {
			close(ch)
		}
	}

	outcome.Fruits = fruits

	return outcome
}

func (testMatter *TestMatter) exec(testTarget *TestTarget, testCase *TestCase) *Detail {

	detail := testTarget.Exec(testMatter.Client, testCase)

	// if result not set (this means execution was successful).
	if detail.Result == "" {

		if testMatter.TestMethod != nil {
			// use custom testing method.
			res, ers := testMatter.TestMethod.Exec(testMatter.Client, testCase, detail)
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

	detail.IsExpected = detail.Result == testTarget.Expect

	return detail

}

func (testMatter *TestMatter) goEach(delegateFunc func(int, int, *TestTarget, *TestCase)) {
	for i, testTarget := range testMatter.TestTargets {
		for j, testCase := range testMatter.TestCases {
			go delegateFunc(i, j, testTarget, testCase)
		}
	}
}
