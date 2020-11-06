package tester

import (
	"bytes"
	"fmt"
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
	if len(testTargets) == 0 {
		return nil
	}

	// テストケース
	testCases := MakeTestCases(basepath, config.TestCase)
	if len(testCases) == 0 {
		return nil
	}

	// テストメソッド
	testMethod := NewTestMethod(basepath, config.TestMethod)

	return &TestMatter{dirname, client, testMethod, testTargets, testCases}
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

	for i, testTarget := range testMatter.TestTargets {
		for j, testCase := range testMatter.TestCases {
			go func(i, j int, testTarget *TestTarget, testCase *TestCase) {
				fruits[i].Details[j] = testMatter.exec(testTarget, testCase)
				ch <- i
			}(i, j, testTarget, testCase)
		}
	}

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
	result, time, output, e := testMatter.do(testTarget.Language, testTarget.SourceCode, testCase.Input)

	if result == "" {
		// input for test_method goes in this format.
		// output + "\0" + input + "\0" + answer
		input := strings.Join([]string{output, testCase.Input.String(), testCase.Answer.String()}, "\000")
		if testMatter.TestMethod != nil {
			res, _, out, ers := testMatter.do(testMatter.TestMethod.Language, testMatter.TestMethod.SourceCode, bytes.NewBufferString(input))
			if res == "" {
				result = strings.TrimRight(out, "\n")
				e += ers
			} else {
				result = fmt.Sprintf("METHOD %s", res)
				e += ers
			}
		} else {
			if output == testCase.Answer.String() {
				result = "PASS"
			} else {
				result = "FAIL"
			}
		}
	}

	isExpected := result == testTarget.Expect

	return &Detail{testCase.Name, result, isExpected, time, output, e}
}

func (testMatter *TestMatter) do(language string, sourceCode, input *bytes.Buffer) (string, string, string, string) {

	res1, err := testMatter.Client.RunnersCreate(language, sourceCode, input)
	if err != nil {
		if err.Code >= 500 {
			return "SERVER ERROR", "", "", err.Error()
		} else if err.Code >= 400 {
			return "CLIENT ERROR", "", "", err.Error()
		} else {
			return "TESTER ERROR", "", "", err.Error()
		}
	}

	res2, err := testMatter.Client.RunnersGetDetails(res1.ID)
	if err != nil {
		if err.Code >= 500 {
			return "SERVER ERROR", "", "", err.Error()
		} else if err.Code >= 400 {
			return "CLIENT ERROR", "", "", err.Error()
		} else {
			return "TESTER ERROR", "", "", err.Error()
		}
	}

	if !(res2.BuildResult == "" || res2.BuildResult == "success") {
		return fmt.Sprintf("BUILD %s", strings.ToUpper(res2.BuildResult)), "", "", res2.BuildSTDERR
	}

	if res2.Result != "success" {
		return fmt.Sprintf("EXECUTION %s", strings.ToUpper(res2.Result)), "", "", res2.STDERR
	}

	return "", res2.Time, res2.STDOUT, res2.STDERR
}
