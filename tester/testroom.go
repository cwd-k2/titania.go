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

type execinfo struct {
	UnitName string
	Info     *TestInfo
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
		baseDirectoryPath,
		languages,
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

func (testRoom *TestRoom) Exec() []*TestOver {
	ch := make(chan execinfo)
	view := InitTestView(testRoom.TestUnits, testRoom.TestCases)
	over := make(map[string]*TestOver)

	for unitName, testUnit := range testRoom.TestUnits {
		testOver := new(TestOver)
		testOver.UnitName = unitName
		testOver.Language = testUnit.Language
		over[unitName] = testOver
	}

	testRoom.goEach(func(testUnit *TestUnit, testCase *TestCase) {
		ch <- testRoom.execTest(testUnit, testCase)
	})

	// 出力する
	view.Start()

	curr := 0
	stop := len(testRoom.TestUnits) * len(testRoom.TestCases)

	for exec := range ch {
		curr++
		view.Refresh(exec.UnitName)

		over[exec.UnitName].Details =
			append(over[exec.UnitName].Details, exec.Info)

		if curr == stop {
			close(ch)
		}
	}

	var results []*TestOver

	for _, testOver := range over {
		sort.Slice(testOver.Details, func(i, j int) bool {
			return testOver.Details[i].CaseName < testOver.Details[j].CaseName
		})
		results = append(results, testOver)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].UnitName < results[j].UnitName
	})

	return results
}

func (testRoom *TestRoom) execTest(
	testUnit *TestUnit, testCase *TestCase) execinfo {

	unitName := testUnit.Name
	caseName := testCase.Name

	testInfo := new(TestInfo)
	testInfo.CaseName = caseName

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := testRoom.Client.Do(testUnit.SourceCode, testUnit.Language, testCase.Input)

	if err != nil {
		if err.Code >= 500 {
			testInfo.Result = "SERVER ERROR"
		} else if err.Code >= 400 {
			testInfo.Result = "CLIENT ERROR"
		} else {
			testInfo.Result = "TESTER ERROR"
		}
		testInfo.Error = err.Error()
		testInfo.Time = ""
		return execinfo{unitName, testInfo}
	}

	// ビルドエラー
	if !(resp.BuildResult == "success" ||
		resp.BuildResult == "") {
		testInfo.Result = fmt.Sprintf("BUILD %s", strings.ToUpper(resp.BuildResult))
		testInfo.Error = resp.BuildSTDERR
		testInfo.Time = ""
		return execinfo{unitName, testInfo}
	}

	// 実行時エラー
	if resp.Result != "success" {
		testInfo.Result = fmt.Sprintf("EXECUTION %s", strings.ToUpper(resp.Result))
		testInfo.Error = resp.STDERR
		testInfo.Time = ""
		return execinfo{unitName, testInfo}
	}

	// 出力が正しいかどうか
	if resp.STDOUT == testCase.Output {
		testInfo.Result = "PASS"
	} else {
		testInfo.Result = "FAIL"
	}

	testInfo.Time = resp.Time
	return execinfo{unitName, testInfo}

}

func (testRoom *TestRoom) goEach(delegateFunc func(*TestUnit, *TestCase)) {
	for _, testUnit := range testRoom.TestUnits {
		for _, testCase := range testRoom.TestCases {
			go delegateFunc(testUnit, testCase)
		}
	}
}
