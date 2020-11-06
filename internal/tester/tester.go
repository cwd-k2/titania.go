package tester

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cwd-k2/titania.go/internal/client"
)

// Tester
// contains paiza.io API client, SourceCodes, and TestCases
// physically, this stands for a directory.
type Tester struct {
	Name        string
	Client      *client.Client
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
}

// Newtester
// returns *tester
func NewTester(dirname string, languages []string) *Tester {
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

	return &Tester{dirname, client, testMethod, testTargets, testCases}
}

func MakeTesters(directories, languages []string) []*Tester {
	testers := make([]*Tester, 0, len(directories))
	for _, dirname := range directories {
		tester := NewTester(dirname, languages)
		if tester != nil {
			testers = append(testers, tester)
		}
	}
	return testers
}

func (tester *Tester) Exec(view View) *Outcome {
	curr := 0
	stop := len(tester.TestTargets) * len(tester.TestCases)

	ch := make(chan int, stop)
	fruits := make([]*Fruit, len(tester.TestTargets))

	outcome := new(Outcome)
	outcome.Name = tester.Name
	if tester.TestMethod != nil {
		outcome.TestMethod = tester.TestMethod.Name
	} else {
		outcome.TestMethod = "default"
	}

	for i, testTarget := range tester.TestTargets {
		fruit := new(Fruit)
		fruit.TestTarget = testTarget.Name
		fruit.Language = testTarget.Language
		fruit.Expect = testTarget.Expect
		fruit.Details = make([]*Detail, len(tester.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	for i, testTarget := range tester.TestTargets {
		for j, testCase := range tester.TestCases {
			go func(i, j int, testTarget *TestTarget, testCase *TestCase) {
				fruits[i].Details[j] = tester.exec(testTarget, testCase)
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

func (tester *Tester) exec(testTarget *TestTarget, testCase *TestCase) *Detail {
	result, time, output, e := tester.do(testTarget.Language, testTarget.SourceCode, testCase.Input)

	if result == "" {
		// input for test_method goes in this format.
		// output + "\0" + input + "\0" + answer
		input := strings.Join([]string{output, testCase.Input.String(), testCase.Answer.String()}, "\000")
		if tester.TestMethod != nil {
			res, _, out, ers := tester.do(tester.TestMethod.Language, tester.TestMethod.SourceCode, bytes.NewBufferString(input))
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

func (tester *Tester) do(language string, sourceCode, input *bytes.Buffer) (string, string, string, string) {

	res1, err := tester.Client.RunnersCreate(language, sourceCode, input)
	if err != nil {
		if err.Code >= 500 {
			return "SERVER ERROR", "", "", err.Error()
		} else if err.Code >= 400 {
			return "CLIENT ERROR", "", "", err.Error()
		} else {
			return "TESTER ERROR", "", "", err.Error()
		}
	}

	res2, err := tester.Client.RunnersGetDetails(res1.ID)
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

func Exec(directories, languages []string, async bool) []*Outcome {
	testers := MakeTesters(directories, languages)

	if len(testers) == 0 {
		return nil
	}

	outcomes := make([]*Outcome, len(testers))

	if async {

		wg := new(sync.WaitGroup)

		for i, tester := range testers {
			wg.Add(1)

			go func(i int, tester *Tester) {
				defer wg.Done()
				view := InitView(tester, true)
				outcome := tester.Exec(view)
				outcomes[i] = outcome
			}(i, tester)

		}

		wg.Wait()

	} else {

		for i, tester := range testers {
			view := InitView(tester, false)
			outcome := tester.Exec(view)
			outcomes[i] = outcome
		}

	}

	return outcomes
}
