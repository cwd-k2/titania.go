package tester

import (
	"fmt"
	"path/filepath"
	"strings"

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
func NewTestUnit(dirname string, languages []string) *TestUnit {
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
	targets := MakeTestTargets(basepath, languages, config.TestTarget)
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

	t.view.Draw()

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

	outcome := &Outcome{t.Name, "default", fruits}
	if t.TestMethod != nil {
		outcome.TestMethod = t.TestMethod.Name
	}

	return outcome
}

func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *Detail {
	// TODO: refactoring
	result, time, stdout, stderr := t.do(target.Language, target.SourceCode, tcase.Input)

	if result == "" {
		// TODO: Method Execution `on` specified result.
		if t.TestMethod != nil {
			// input for test_method goes in this format.
			// output + "\0" + input + "\0" + answer
			input := strings.Join([]string{stdout, tcase.Input, tcase.Answer}, "\000")

			res, _, out, ers := t.do(t.TestMethod.Language, t.TestMethod.SourceCode, input)

			if res == "" {
				result = strings.TrimRight(out, "\n")
				stderr += ers
			} else {
				result = fmt.Sprintf("METHOD %s", res)
				stderr += ers
			}

		} else {
			// simple comparison
			if stdout == tcase.Answer {
				result = "PASS"
			} else {
				result = "FAIL"
			}

		}
	}

	isExpected := result == target.Expect

	return &Detail{tcase.Name, result, isExpected, time, stdout, stderr}
}

func (t *TestUnit) do(language string, sourceCode, input string) (string, string, string, string) {
	// TODO: refactoring

	req1 := &paizaio.RunnersCreateRequest{
		Language:        language,
		SourceCode:      sourceCode,
		Input:           input,
		Longpoll:        true,
		LongpollTimeout: 16,
	}

	res1, err := t.Client.RunnersCreate(req1)
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

	req2 := &paizaio.RunnersGetDetailsRequest{
		ID: res1.ID,
	}

	res2, err := t.Client.RunnersGetDetails(req2)
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

	if res2.BuildExitCode != 0 {
		result := fmt.Sprintf("BUILD %s", strings.ToUpper(res2.BuildResult))
		return result, res2.BuildTime, res2.BuildSTDOUT, res2.BuildSTDERR
	}

	if res2.ExitCode != 0 || res2.Result != "success" {
		result := fmt.Sprintf("EXECUTION %s", strings.ToUpper(res2.Result))
		return result, res2.Time, res2.STDOUT, res2.STDERR
	}

	return "", res2.Time, res2.STDOUT, res2.STDERR
}
