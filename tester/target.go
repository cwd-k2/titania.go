package tester

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/cwd-k2/titania.go/client"
)

// Target
// contains paiza.io API client, config, SourceCodes, and TestCases
type Target struct {
	Name        string
	Client      *client.Client
	Config      *Config
	Method      *Method
	SourceCodes []*SourceCode
	TestCases   []*TestCase
}

// NewTarget
// returns *Target
func NewTarget(dirname string, languages []string) *Target {
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
	method := NewMethod(basepath, config.TestMethodFileName)

	target := new(Target)
	target.Name = dirname
	target.Client = client
	target.Config = config
	target.SourceCodes = sourceCodes
	target.TestCases = testCases
	target.Method = method

	return target
}

func MakeTargets(directories, languages []string) []*Target {
	targets := make([]*Target, 0, len(directories))
	for _, dirname := range directories {
		target := NewTarget(dirname, languages)
		if target != nil {
			targets = append(targets, target)
		}
	}
	return targets
}

func (target *Target) Exec(view View) *Outcome {
	wg := new(sync.WaitGroup)
	fruits := make([]*Fruit, len(target.SourceCodes))

	for i, sourceCode := range target.SourceCodes {
		fruit := new(Fruit)
		fruit.SourceCode = sourceCode.Name
		fruit.Language = sourceCode.Language
		fruit.Details = make([]*Detail, len(target.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	target.goEachWithWg(wg, func(i, j int, sourceCode *SourceCode, testCase *TestCase) {
		defer wg.Done()
		detail := target.exec(sourceCode, testCase)
		fruits[i].Details[j] = detail
		view.Update(i)
	})

	wg.Wait()

	outcome := new(Outcome)
	outcome.Target = target.Name
	if target.Method != nil {
		outcome.Method = target.Method.Name
	} else {
		outcome.Method = "default"
	}
	outcome.Fruits = fruits

	return outcome
}

func (target *Target) exec(sourceCode *SourceCode, testCase *TestCase) *Detail {

	detail := sourceCode.Exec(target.Client, testCase)

	// if result not set (this means execution was successful).
	if detail.Result == "" {

		if target.Method != nil {
			// use custom testing method.
			res, ers := target.Method.Exec(target.Client, testCase, detail)
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

func (target *Target) goEachWithWg(wg *sync.WaitGroup, delegateFunc func(int, int, *SourceCode, *TestCase)) {
	for i, sourceCode := range target.SourceCodes {
		for j, testCase := range target.TestCases {
			wg.Add(1)
			go delegateFunc(i, j, sourceCode, testCase)
		}
	}
}
