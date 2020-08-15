package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
	"github.com/cwd-k2/titania.go/pretty"
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

// returns map[string]*TestRoom
func MakeTestRooms(directories []string) map[string]*TestRoom {
	testRooms := make(map[string]*TestRoom)
	// この前にディレクトリ直下に titania.json がいるか確認したい
	for _, dirname := range directories {
		// ディレクトリ直下の titania.json を読んで設定を作る

		baseDirectoryPath, err := filepath.Abs(dirname)
		// ここのエラーは公式のドキュメント見てもわからんのだけど何？
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			continue
		}

		configFileName := path.Join(baseDirectoryPath, "titania.json")
		configRawData, err := ioutil.ReadFile(configFileName)
		// File Read 失敗
		if err != nil {
			fmt.Fprintf(os.Stderr, "[SKIP] Couldn't read %s.\n", configFileName)
			continue
		}

		// ようやく設定の構造体を作れる
		config := new(Config)

		// JSON パース失敗
		if err := json.Unmarshal(configRawData, config); err != nil {
			fmt.Fprintf(os.Stderr, "[SKIP] Couldn't parse %s.\n%s\n", configFileName, err)
			continue
		}

		// paiza.io API クライアント
		client := new(client.Client)
		client.Host = config.Host
		client.APIKey = config.APIKey

		// テストケース
		testCases := MakeTestCases(
			baseDirectoryPath,
			config.TestCaseDirectories,
			config.TestCaseInputExtension,
			config.TestCaseOutputExtension)

		// テストユニット
		testUnits := MakeTestUnits(
			baseDirectoryPath,
			config.SourceCodeDirectories)

		testRooms[dirname] = new(TestRoom)
		testRooms[dirname].Client = client
		testRooms[dirname].Config = config
		testRooms[dirname].TestUnits = testUnits
		testRooms[dirname].TestCases = testCases

	}

	return testRooms
}

func (testRoom *TestRoom) Exec() {
	ch := make(chan []string)

	exec := func(unitName string, caseName string) {
		client := testRoom.Client
		testUnit := testRoom.TestUnits[unitName]
		testCase := testRoom.TestCases[caseName]

		// 実際に paiza.io の API を利用して実行結果をもらう
		// この辺も分割したい
		runnersCreateResponse, err :=
			client.RunnersCreate(
				testUnit.SourceCode,
				testUnit.Language,
				testCase.Input)

		if err != nil {
			ch <- []string{unitName, caseName, "TEST FAIL", err.Error()}
			return
		}

		runnersGetDetailsResponse, err :=
			client.RunnersGetDetails(runnersCreateResponse.ID)

		if err != nil {
			ch <- []string{unitName, caseName, "TEST FAIL", err.Error()}
			return
		}

		// この辺も分割したい
		// 返ってきた結果をまとめる

		// ビルドエラー
		if !(runnersGetDetailsResponse.BuildResult == "success" ||
			runnersGetDetailsResponse.BuildResult == "") {
			ch <- []string{
				unitName,
				caseName,
				"BUILD " + strings.ToUpper(runnersGetDetailsResponse.BuildResult),
				runnersGetDetailsResponse.BuildSTDERR,
			}
			return
		}

		// 実行時エラー
		if runnersGetDetailsResponse.Result != "success" {
			ch <- []string{
				unitName,
				caseName,
				strings.ToUpper(runnersGetDetailsResponse.Result),
				runnersGetDetailsResponse.STDERR,
			}
			return
		}

		// 出力が正しいかどうか
		if runnersGetDetailsResponse.STDOUT == testCase.Output {
			ch <- []string{unitName, caseName, "PASS"}
		} else {
			ch <- []string{unitName, caseName, "WRONG ANSWER"}
		}

	}

	testRoom.goEach(exec)

	// 出力する
	// ここをいい感じにしたい
	// 今はとりあえず結果だけを出してるけど
	// SUMMARY と DETAIL を出したいね
	// あと色

	i := 0
	j := len(testRoom.TestCases) * len(testRoom.TestUnits)
	k := len(testRoom.TestUnits)

	testShow := make(map[string][]string)
	var keys []string

	for unitName := range testRoom.TestUnits {
		testShow[unitName] = []string{
			"initializing",
		}
		keys = append(keys, unitName)
	}
	for _, key := range keys {
		fmt.Printf("  [UNIT] %s\n", key)
		fmt.Printf("    [WAIT] %s\n", testShow[key][0])
	}

	for msg := range ch {
		i++
		testShow[msg[0]] = msg[1:]
		pretty.Up(2 * k)
		for _, key := range keys {
			pretty.Erase()
			fmt.Printf("  [UNIT] %s\n", key)
			pretty.Erase()
			fmt.Printf("    [WAIT] %s\n", testShow[key][0])
		}
		if i == j {
			close(ch)
		}
	}

}

func (testRoom *TestRoom) goEach(delegateFunc func(string, string)) {
	for unitName := range testRoom.TestUnits {
		for caseName := range testRoom.TestCases {
			go delegateFunc(unitName, caseName)
		}
	}
}
