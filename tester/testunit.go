package tester

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cwd-k2/titania.go/client"
)

// TestUnit
// contains paiza.io API client, config, TestCodes, and TestCases
type TestUnit struct {
	Name      string
	Client    *client.Client
	Config    *Config
	TestCodes []*TestCode
	TestCases []*TestCase
}

// NewtestUnit
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
	testCodes := MakeTestCodes(basepath, languages, config.SourceCodeDirectories)

	// ソースコードがなければ実行しない
	if len(testCodes) == 0 {
		return nil
	}

	// テストケース
	testCases := MakeTestCases(
		basepath,
		config.TestCaseDirectories,
		config.TestCaseInputExtension,
		config.TestCaseOutputExtension)

	// テストケースがなければ実行しない
	if len(testCases) == 0 {
		return nil
	}

	testUnit := new(TestUnit)
	testUnit.Name = dirname
	testUnit.Client = client
	testUnit.Config = config
	testUnit.TestCodes = testCodes
	testUnit.TestCases = testCases

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

func (testUnit *TestUnit) Exec(view View) []*ShowCode {
	wg := new(sync.WaitGroup)
	fruits := make([]*ShowCode, len(testUnit.TestCodes))

	for i, testCode := range testUnit.TestCodes {
		fruit := new(ShowCode)
		fruit.Name = testCode.Name
		fruit.Language = testCode.Language
		fruit.Details = make([]*ShowCase, len(testUnit.TestCases))
		fruits[i] = fruit
	}

	view.Draw()

	testUnit.goEachWithWg(wg, func(i, j int, testCode *TestCode, testCase *TestCase) {
		defer wg.Done()
		detail := testUnit.exec(testCode, testCase)
		fruits[i].Details[j] = detail
		view.Update(i)
	})

	wg.Wait()

	return fruits
}

func (testUnit *TestUnit) exec(testCode *TestCode, testCase *TestCase) *ShowCase {

	detail := new(ShowCase)
	detail.Name = testCase.Name

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := testUnit.Client.Do(testCode.SourceCode, testCode.Language, testCase.Input)

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
	if !(resp.BuildResult == "success" ||
		resp.BuildResult == "") {
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
	if resp.STDOUT == testCase.Output {
		detail.Result = "PASS"
	} else {
		detail.Result = "FAIL"
	}

	detail.Time = resp.Time
	detail.OutPut = resp.STDOUT
	return detail

}

func (testUnit *TestUnit) goEachWithWg(wg *sync.WaitGroup, delegateFunc func(int, int, *TestCode, *TestCase)) {
	for i, testCode := range testUnit.TestCodes {
		for j, testCase := range testUnit.TestCases {
			wg.Add(1)
			go delegateFunc(i, j, testCode, testCase)
		}
	}
}
