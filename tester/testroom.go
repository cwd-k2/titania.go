package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// Config
// test configs
type Config struct {
	Host                    string   `json:"host"`
	APIKey                  string   `json:"api_key"`
	SourceCodeDirectories   []string `json:"source_code_directories"`
	TestCaseDirectories     []string `json:"test_case_directories"`
	TestCaseInputExtension  string   `json:"test_case_input_extension"`
	TestCaseOutputExtension string   `json:"test_case_output_extension"`
	MaxProcesses            uint     `json:"max_processes"`
}

// TestRoom
// contains paiza.io API client, config, and map of TestUnits, map of TestCases
type TestRoom struct {
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

	// ディレクトリ直下に titania.json がいるか確認したい
	configFileName := path.Join(baseDirectoryPath, "titania.json")
	if match, _ := filepath.Glob(configFileName); len(match) == 0 {
		return nil
	}

	// ディレクトリ直下の titania.json を読んで設定を作る
	configRawData, err := ioutil.ReadFile(configFileName)

	// File Read 失敗
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read %s.\n", configFileName)
		return nil
	}

	// ようやく設定の構造体を作れる
	config := new(Config)

	// JSON パース失敗
	if err := json.Unmarshal(configRawData, config); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Couldn't parse %s.\n%s\n",
			configFileName, err)
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

func (testRoom *TestRoom) Exec() []*TestInfo {
	ch := make(chan *TestInfo)
	view := InitTestView(testRoom.TestUnits, testRoom.TestCases)
	var results []*TestInfo

	testRoom.goEach(func(unitName string, caseName string) {
		ch <- testRoom.execTest(unitName, caseName)
	})

	// 出力する
	view.Start()

	i := 0
	j := len(testRoom.TestUnits) * len(testRoom.TestCases)

	for testInfo := range ch {
		i++
		view.Refresh(testInfo)
		results = append(results, testInfo)
		if i == j {
			close(ch)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].UnitName == results[j].UnitName {
			return results[i].CaseName < results[j].CaseName
		} else {
			return results[i].UnitName < results[j].UnitName
		}
	})
	return results
}

func (testRoom *TestRoom) execTest(unitName string, caseName string) *TestInfo {
	client := testRoom.Client
	testUnit := testRoom.TestUnits[unitName]
	testCase := testRoom.TestCases[caseName]
	testInfo := new(TestInfo)
	testInfo.UnitName = unitName
	testInfo.CaseName = caseName
	testInfo.Language = strings.ToUpper(testUnit.Language)

	// 実際に paiza.io の API を利用して実行結果をもらう
	// この辺も分割したい
	runnersCreateResponse, err :=
		client.RunnersCreate(
			testUnit.SourceCode,
			testUnit.Language,
			testCase.Input)

	if err != nil {
		testInfo.Result = "TEST FAIL"
		testInfo.Error = err.Error()
		testInfo.Time = ""
		return testInfo
	}

	runnersGetDetailsResponse, err :=
		client.RunnersGetDetails(runnersCreateResponse.ID)

	if err != nil {
		testInfo.Result = "TEST FAIL"
		testInfo.Error = err.Error()
		testInfo.Time = ""
		return testInfo
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
		return testInfo
	}

	// 実行時エラー
	if runnersGetDetailsResponse.Result != "success" {
		testInfo.Result = strings.ToUpper(runnersGetDetailsResponse.Result)
		testInfo.Error = runnersGetDetailsResponse.STDERR
		testInfo.Time = ""
		return testInfo
	}

	// 出力が正しいかどうか
	if runnersGetDetailsResponse.STDOUT == testCase.Output {
		testInfo.Result = "PASS"
		testInfo.Time = runnersGetDetailsResponse.Time
		return testInfo
	} else {
		testInfo.Result = "WRONG ANSWER"
		testInfo.Time = runnersGetDetailsResponse.Time
		return testInfo
	}

}

func (testRoom *TestRoom) goEach(delegateFunc func(string, string)) {
	for unitName := range testRoom.TestUnits {
		for caseName := range testRoom.TestCases {
			go delegateFunc(unitName, caseName)
		}
	}
}
