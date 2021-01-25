package tester

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/paizaio"
	"github.com/cwd-k2/titania.go/pkg/viewer"
)

// TestUnit
// contains paiza.io API client, SourceCodes, and TestCases
// physically, this stands for a directory.
type TestUnit struct {
	Name        string
	Client      *paizaio.Client
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
}

// NewTestUnit
// returns *tester
func NewTestUnit(dirname string, languages []string) *TestUnit {
	basepath, err := filepath.Abs(dirname)
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
	client := paizaio.NewClient(config.ClientConfig)

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

	return &TestUnit{dirname, client, testMethod, testTargets, testCases}
}

func MakeTestUnits(directories, languages []string) []*TestUnit {
	ts := make([]*TestUnit, 0, len(directories))
	for _, dirname := range directories {
		t := NewTestUnit(dirname, languages)
		if t != nil {
			ts = append(ts, t)
		}
	}
	return ts
}

func (t *TestUnit) Exec(view viewer.Viewer) *Outcome {
	curr := 0
	stop := len(t.TestTargets) * len(t.TestCases)

	ch := make(chan int, stop)
	fruits := make([]*Fruit, len(t.TestTargets))

	outcome := new(Outcome)
	outcome.Name = t.Name
	if t.TestMethod != nil {
		outcome.TestMethod = t.TestMethod.Name
	} else {
		outcome.TestMethod = "default"
	}

	for i, testTarget := range t.TestTargets {
		fruits[i] = &Fruit{
			testTarget.Name,
			testTarget.Language,
			testTarget.Expect,
			make([]*Detail, len(t.TestCases)),
		}
	}

	view.Draw()

	for i, testTarget := range t.TestTargets {
		for j, testCase := range t.TestCases {
			go func(i, j int, testTarget *TestTarget, testCase *TestCase) {
				fruits[i].Details[j] = t.exec(testTarget, testCase)
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

func (t *TestUnit) exec(testTarget *TestTarget, testCase *TestCase) *Detail {
	result, time, output, e := t.do(testTarget.Language, testTarget.SourceCode, testCase.Input)

	if result == "" {
		// input for test_method goes in this format.
		// output + "\0" + input + "\0" + answer
		input := strings.Join([]string{output, testCase.Input, testCase.Answer}, "\000")

		if t.TestMethod != nil {
			res, _, out, ers := t.do(t.TestMethod.Language, t.TestMethod.SourceCode, input)

			if res == "" {
				result = strings.TrimRight(out, "\n")
				e += ers
			} else {
				result = fmt.Sprintf("METHOD %s", res)
				e += ers
			}

		} else {

			if output == testCase.Answer {
				result = "PASS"
			} else {
				result = "FAIL"
			}

		}
	}

	isExpected := result == testTarget.Expect

	return &Detail{testCase.Name, result, isExpected, time, output, e}
}

func (t *TestUnit) do(language string, sourceCode, input string) (string, string, string, string) {

	res1, err := t.Client.RunnersCreate(language, sourceCode, input)
	if err != nil {
		switch err := err.(type) {
		case paizaio.ServerError:
			return "SERVER ERROR", "", "", err.Error()
		case paizaio.ClientError:
			return "CLIENT ERROR", "", "", err.Error()
		default:
			return "TESTER ERROR", "", "", err.Error()
		}
	}

	res2, err := t.Client.RunnersGetDetails(res1.ID)
	if err != nil {
		switch err := err.(type) {
		case paizaio.ServerError:
			return "SERVER ERROR", "", "", err.Error()
		case paizaio.ClientError:
			return "CLIENT ERROR", "", "", err.Error()
		default:
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
