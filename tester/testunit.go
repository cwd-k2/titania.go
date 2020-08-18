package tester

import (
	"fmt"
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

	target := new(Target)
	target.Name = dirname
	target.Client = client
	target.Config = config
	target.SourceCodes = sourceCodes
	target.TestCases = testCases

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
	outcome.Fruits = fruits

	return outcome
}

func (target *Target) exec(sourceCode *SourceCode, testCase *TestCase) *Detail {

	detail := new(Detail)
	detail.TestCase = testCase.Name

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := target.Client.Do(sourceCode.SourceCode, sourceCode.Language, testCase.Input)

	if err != nil {
		if err.Code >= 500 {
			detail.Result = "SERVER ERROR"
		} else if err.Code >= 400 {
			detail.Result = "CLIENT ERROR"
		} else {
			detail.Result = "TESTER ERROR"
		}
		detail.Error = err.Error()
		return detail
	}

	// ビルドエラー
	if !(resp.BuildResult == "success" || resp.BuildResult == "") {
		detail.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(resp.BuildResult))
		detail.Error = resp.BuildSTDERR
		return detail
	}

	// 実行時エラー
	if resp.Result != "success" {
		detail.Result = fmt.Sprintf("EXECUTION %s", strings.ToUpper(resp.Result))
		detail.Error = resp.STDERR
		return detail
	}

	// 出力が正しいかどうか
	if resp.STDOUT == testCase.Answer {
		detail.Result = "PASS"
	} else {
		detail.Result = "FAIL"
	}

	detail.Time = resp.Time
	detail.OutPut = resp.STDOUT
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
