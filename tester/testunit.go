package tester

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// TestUnit
// contains paiza.io API client, config, and map of Codes, map of TestCases
type TestUnit struct {
	Name      string
	Client    *client.Client
	Config    *Config
	TestCodes map[string]*TestCode
	TestCases map[string]*TestCase
}

// NewtestUnit
// returns *TestUnit
func NewTestUnit(dirname string, languages []string) *TestUnit {
	baseDirectoryPath, err := filepath.Abs(dirname)
	// ここのエラーは公式のドキュメント見てもわからんのだけど何？
	if err != nil {
		println(err)
		return nil
	}

	// 設定
	config := NewConfig(baseDirectoryPath)
	if config == nil {
		return nil
	}

	// paiza.io API クライアント
	client := new(client.Client)
	client.Host = config.Host
	client.APIKey = config.APIKey

	// ソースコード
	testCodes := MakeTestCodes(
		baseDirectoryPath, languages,
		config.SourceCodeDirectories)

	// ソースコードがなければ実行しない
	if len(testCodes) == 0 {
		return nil
	}

	// テストケース
	testCases := MakeTestCases(
		baseDirectoryPath,
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

func (testUnit *TestUnit) Exec() []*ShowCode {
	ch := make(chan string)

	view := InitView(testUnit.TestCodes, testUnit.TestCases)
	view.Draw(testUnit.Name)

	overs := make(map[string]*ShowCode)

	for name, testCode := range testUnit.TestCodes {
		over := new(ShowCode)
		over.Name = name
		over.Language = testCode.Language
		overs[name] = over
	}

	testUnit.goEach(func(testCodes *TestCode, testCase *TestCase) {
		unitName, detail := testUnit.execTest(testCodes, testCase)
		overs[unitName].Details = append(overs[unitName].Details, detail)

		ch <- unitName
	})

	curr := 0
	stop := len(testUnit.TestCodes) * len(testUnit.TestCases)

	for unitName := range ch {
		curr++
		view.Update(unitName)

		if curr == stop {
			close(ch)
		}
	}

	var fruits []*ShowCode

	for _, over := range overs {
		sort.Slice(over.Details, func(i, j int) bool {
			return over.Details[i].Name < over.Details[j].Name
		})
		fruits = append(fruits, over)
	}

	sort.Slice(fruits, func(i, j int) bool {
		return fruits[i].Name < fruits[j].Name
	})

	return fruits
}

func (testUnit *TestUnit) execTest(
	testCode *TestCode, testCase *TestCase) (string, *ShowCase) {

	unitName := testCode.Name
	caseName := testCase.Name

	showCase := new(ShowCase)
	showCase.Name = caseName

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := testUnit.Client.Do(testCode.SourceCode, testCode.Language, testCase.Input)

	if err != nil {
		if err.Code >= 500 {
			showCase.Result = "SERVER ERROR"
		} else if err.Code >= 400 {
			showCase.Result = "CLIENT ERROR"
		} else {
			showCase.Result = "TESTER ERROR"
		}
		showCase.Error = err.Error()
		return unitName, showCase
	}

	// ビルドエラー
	if !(resp.BuildResult == "success" ||
		resp.BuildResult == "") {
		showCase.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(resp.BuildResult))
		showCase.Error = resp.BuildSTDERR
		return unitName, showCase
	}

	// 実行時エラー
	if resp.Result != "success" {
		showCase.Result = fmt.Sprintf("EXECUTION %s", strings.ToUpper(resp.Result))
		showCase.Error = resp.STDERR
		return unitName, showCase
	}

	// 出力が正しいかどうか
	if resp.STDOUT == testCase.Output {
		showCase.Result = "PASS"
	} else {
		showCase.Result = "FAIL"
	}

	showCase.Time = resp.Time
	showCase.OutPut = resp.STDOUT
	return unitName, showCase

}

func (testUnit *TestUnit) goEach(delegateFunc func(*TestCode, *TestCase)) {
	for _, testCode := range testUnit.TestCodes {
		for _, testCase := range testUnit.TestCases {
			go delegateFunc(testCode, testCase)
		}
	}
}
