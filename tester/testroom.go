package tester

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// TestRoom
// contains paiza.io API client, config, and map of TestUnits, map of TestCases
type TestRoom struct {
	Name      string
	Client    *client.Client
	Config    *Config
	TestUnits map[string]*TestUnit
	TestCases map[string]*TestCase
}

// NewTestRoom
// returns *TestRoom
func NewTestRoom(dirname string, languages []string) *TestRoom {
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

	// テストユニット
	testUnits := MakeTestUnits(
		baseDirectoryPath, languages,
		config.SourceCodeDirectories)

	// テストユニットがなければ実行しない
	if len(testUnits) == 0 {
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

	testRoom := new(TestRoom)
	testRoom.Name = dirname
	testRoom.Client = client
	testRoom.Config = config
	testRoom.TestUnits = testUnits
	testRoom.TestCases = testCases

	return testRoom
}

func (testRoom *TestRoom) Exec() []*ShowCode {
	ch := make(chan string)

	view := InitView(testRoom.TestUnits, testRoom.TestCases)
	view.Draw()

	overs := make(map[string]*ShowCode)

	for unitName, testUnit := range testRoom.TestUnits {
		over := new(ShowCode)
		over.Name = unitName
		over.Language = testUnit.Language
		overs[unitName] = over
	}

	testRoom.goEach(func(testUnit *TestUnit, testCase *TestCase) {
		unitName, detail := testRoom.execTest(testUnit, testCase)
		overs[unitName].Details = append(overs[unitName].Details, detail)

		ch <- unitName
	})

	curr := 0
	stop := len(testRoom.TestUnits) * len(testRoom.TestCases)

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

func (testRoom *TestRoom) execTest(
	testUnit *TestUnit, testCase *TestCase) (string, *ShowCase) {

	unitName := testUnit.Name
	caseName := testCase.Name

	showCase := new(ShowCase)
	showCase.Name = caseName

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := testRoom.Client.Do(testUnit.SourceCode, testUnit.Language, testCase.Input)

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

func (testRoom *TestRoom) goEach(delegateFunc func(*TestUnit, *TestCase)) {
	for _, testUnit := range testRoom.TestUnits {
		for _, testCase := range testRoom.TestCases {
			go delegateFunc(testUnit, testCase)
		}
	}
}
