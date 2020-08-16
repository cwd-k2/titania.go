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
		over[unitName] = new(TestOver)
		over[unitName].UnitName = unitName
		over[unitName].Language = testUnit.Language
	}

	testRoom.goEach(func(unitName string, caseName string) {
		ch <- testRoom.execTest(unitName, caseName)
	})

	// 出力する
	view.Start()

	i := 0
	j := len(testRoom.TestUnits) * len(testRoom.TestCases)

	for exei := range ch {
		i++
		view.Refresh(exei.UnitName)

		over[exei.UnitName].Details =
			append(over[exei.UnitName].Details, exei.Info)

		if i == j {
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
	unitName string, caseName string) execinfo {

	testUnit := testRoom.TestUnits[unitName]
	testCase := testRoom.TestCases[caseName]
	testInfo := new(TestInfo)
	testInfo.CaseName = caseName

	// 実際に paiza.io の API を利用して実行結果をもらう
	// この辺も分割したい
	runnersCreateResponse, err :=
		testRoom.Client.RunnersCreate(
			testUnit.SourceCode,
			testUnit.Language,
			testCase.Input)

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

	runnersGetDetailsResponse, err :=
		testRoom.Client.RunnersGetDetails(runnersCreateResponse.ID)

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
	if !(runnersGetDetailsResponse.BuildResult == "success" ||
		runnersGetDetailsResponse.BuildResult == "") {
		testInfo.Result =
			fmt.Sprintf(
				"BUILD %s",
				strings.ToUpper(runnersGetDetailsResponse.BuildResult))
		testInfo.Error = runnersGetDetailsResponse.BuildSTDERR
		testInfo.Time = ""
		return execinfo{unitName, testInfo}
	}

	// 実行時エラー
	if runnersGetDetailsResponse.Result != "success" {
		testInfo.Result =
			fmt.Sprintf(
				"EXECUTION %s",
				strings.ToUpper(runnersGetDetailsResponse.Result))
		testInfo.Error = runnersGetDetailsResponse.STDERR
		testInfo.Time = ""
		return execinfo{unitName, testInfo}
	}

	// 出力が正しいかどうか
	if runnersGetDetailsResponse.STDOUT == testCase.Output {
		testInfo.Result = "PASS"
	} else {
		testInfo.Result = "FAIL"
	}

	testInfo.Time = runnersGetDetailsResponse.Time
	return execinfo{unitName, testInfo}

}

func (testRoom *TestRoom) goEach(delegateFunc func(string, string)) {
	for unitName := range testRoom.TestUnits {
		for caseName := range testRoom.TestCases {
			go delegateFunc(unitName, caseName)
		}
	}
}
