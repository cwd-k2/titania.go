package tester

import (
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// TestTopic
// contains paiza.io API client, SourceCodes, and TestCases
// physically, this stands for a directory.
type TestTopic struct {
	Name        string
	Client      *client.Client
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
}

// NewTestTopic
// returns *TestTopic
func NewTestTopic(dirname string, languages []string) *TestTopic {
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

	testTopic := new(TestTopic)
	testTopic.Name = dirname
	testTopic.Client = client
	testTopic.TestTargets = testTargets
	testTopic.TestCases = testCases
	testTopic.TestMethod = testMethod

	return testTopic
}

func MakeTestTopics(directories, languages []string) []*TestTopic {
	testTopics := make([]*TestTopic, 0, len(directories))
	for _, dirname := range directories {
		testTopic := NewTestTopic(dirname, languages)
		if testTopic != nil {
			testTopics = append(testTopics, testTopic)
		}
	}
	return testTopics
}

func (testTopic *TestTopic) Exec(view View) *Outcome {
	ch := make(chan int)
	fruits := make([]*Fruit, len(testTopic.TestTargets))

	outcome := new(Outcome)
	outcome.TestTopic = testTopic.Name
	if testTopic.TestMethod != nil {
		outcome.TestMethod = testTopic.TestMethod.Name
	} else {
		outcome.TestMethod = "default"
	}

	curr := 0
	stop := len(testTopic.TestTargets) * len(testTopic.TestCases)

	for i, testTarget := range testTopic.TestTargets {
		fruit := new(Fruit)
		fruit.TestTarget = testTarget.Name
		fruit.Language = testTarget.Language
		fruit.Expect = testTarget.Expect
		fruit.Details = make([]*Detail, len(testTopic.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	testTopic.goEach(func(i, j int, testTarget *TestTarget, testCase *TestCase) {
		detail := testTopic.exec(testTarget, testCase)
		fruits[i].Details[j] = detail
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

func (testTopic *TestTopic) exec(testTarget *TestTarget, testCase *TestCase) *Detail {

	detail := testTarget.Exec(testTopic.Client, testCase)

	// if result not set (this means execution was successful).
	if detail.Result == "" {

		if testTopic.TestMethod != nil {
			// use custom testing method.
			res, ers := testTopic.TestMethod.Exec(testTopic.Client, testCase, detail)
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

func (testTopic *TestTopic) goEach(delegateFunc func(int, int, *TestTarget, *TestCase)) {
	for i, testTarget := range testTopic.TestTargets {
		for j, testCase := range testTopic.TestCases {
			go delegateFunc(i, j, testTarget, testCase)
		}
	}
}
