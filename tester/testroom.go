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

func (testRoom *TestRoom) Exec() []*ShowUnit {
	ch := make(chan string)

	view := InitView(testRoom.TestUnits, testRoom.TestCases)
	view.Draw()

	overs := make(map[string]*ShowUnit)

	for unitName, testUnit := range testRoom.TestUnits {
		over := new(ShowUnit)
		over.UnitName = unitName
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

	var fruits []*ShowUnit

	for _, over := range overs {
		sort.Slice(over.Details, func(i, j int) bool {
			return over.Details[i].CaseName < over.Details[j].CaseName
		})
		fruits = append(fruits, over)
	}

	sort.Slice(fruits, func(i, j int) bool {
		return fruits[i].UnitName < fruits[j].UnitName
	})

	return fruits
}

func (testRoom *TestRoom) execTest(
	testUnit *TestUnit, testCase *TestCase) (string, *ShowCase) {

	unitName := testUnit.Name
	caseName := testCase.Name

	ShowCase := new(ShowCase)
	ShowCase.CaseName = caseName

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := testRoom.Client.Do(testUnit.SourceCode, testUnit.Language, testCase.Input)

	if err != nil {
		if err.Code >= 500 {
			ShowCase.Result = "SERVER ERROR"
		} else if err.Code >= 400 {
			ShowCase.Result = "CLIENT ERROR"
		} else {
			ShowCase.Result = "TESTER ERROR"
		}
		ShowCase.Error = err.Error()
		return unitName, ShowCase
	}

	// ビルドエラー
	if !(resp.BuildResult == "success" ||
		resp.BuildResult == "") {
		ShowCase.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(resp.BuildResult))
		ShowCase.Error = resp.BuildSTDERR
		return unitName, ShowCase
	}

	// 実行時エラー
	if resp.Result != "success" {
		ShowCase.Result = fmt.Sprintf("EXECUTION %s", strings.ToUpper(resp.Result))
		ShowCase.Error = resp.STDERR
		return unitName, ShowCase
	}

	// 出力が正しいかどうか
	if resp.STDOUT == testCase.Output {
		ShowCase.Result = "PASS"
	} else {
		ShowCase.Result = "FAIL"
	}

	ShowCase.Time = resp.Time
	ShowCase.OutPut = resp.STDOUT
	return unitName, ShowCase

}

func (testRoom *TestRoom) goEach(delegateFunc func(*TestUnit, *TestCase)) {
	for _, testUnit := range testRoom.TestUnits {
		for _, testCase := range testRoom.TestCases {
			go delegateFunc(testUnit, testCase)
		}
	}
}
