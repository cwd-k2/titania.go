package tester

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/internal/pkg/viewer"
	"github.com/cwd-k2/titania.go/pkg/paizaio"
)

// TestUnit
// physically, this stands for a directory.
type TestUnit struct {
	Name        string
	Client      *paizaio.Client
	TestMethod  *TestMethod
	TestTargets []*TestTarget
	TestCases   []*TestCase
	view        viewer.Viewer
}

// returns *TestUnit
func NewTestUnit(dirname string, languages []string, quiet bool) *TestUnit {
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

	if quiet {
		// Viewer
		view := viewer.NewQuietView(dirname, len(targets)*len(tcases))
		return &TestUnit{dirname, client, tmethod, targets, tcases, view}
	} else {
		// Viewer
		indices := make([]string, 0)
		for _, target := range targets {
			indices = append(indices, target.Name)
		}
		view := viewer.NewFancyView(dirname, len(targets), len(tcases), indices)
		return &TestUnit{dirname, client, tmethod, targets, tcases, view}
	}

}

type detailstruct struct {
	I int
	J int
	D *Detail
}

func (t *TestUnit) Exec() *Outcome {
	curr := 0
	stop := len(t.TestTargets) * len(t.TestCases)

	fruits := make([]*Fruit, 0)
	for _, target := range t.TestTargets {
		fruits = append(fruits, &Fruit{target.Name, target.Language, target.Expect, make([]*Detail, len(t.TestCases))})
	}

	t.view.Draw()

	ch := make(chan *detailstruct, stop)

	for i, target := range t.TestTargets {
		for j, tcase := range t.TestCases {
			go func(i, j int, target *TestTarget, tcase *TestCase) {
				ch <- &detailstruct{i, j, t.exec(target, tcase)}
			}(i, j, target, tcase)
		}
	}

	for d := range ch {
		t.view.Update(d.I)
		fruits[d.I].Details[d.J] = d.D

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

// TODO: refactoring
func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *Detail {
	result, time, output, e := t.do(target.Language, target.SourceCode, tcase.Input)

	if result == "" {
		// input for test_method goes in this format.
		// output + "\0" + input + "\0" + answer
		input := strings.Join([]string{output, tcase.Input, tcase.Answer}, "\000")

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

			if output == tcase.Answer {
				result = "PASS"
			} else {
				result = "FAIL"
			}

		}
	}

	isExpected := result == target.Expect

	return &Detail{tcase.Name, result, isExpected, time, output, e}
}

// TODO: refactoring
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
