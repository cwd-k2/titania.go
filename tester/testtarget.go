package tester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// TestTarget
// contains source code, its language
type TestTarget struct {
	Name       string
	Language   string
	SourceCode string
	Expect     string
}

// returns []*TestTarget
func MakeTestTargets(
	basepath string,
	languages []string,
	testTartgetDirectories []string) []*TestTarget {

	tmp0 := make([][]*TestTarget, 0, len(testTartgetDirectories))
	length := 0

	for _, dirname := range testTartgetDirectories {
		// ソースファイル
		pattern := filepath.Join(basepath, dirname, "*.*")
		filenames, err := filepath.Glob(pattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err.Error())
			continue
		}

		tmp1 := make([]*TestTarget, 0, len(filenames))

		for _, filename := range filenames {
			name := strings.Replace(filename, basepath+string(filepath.Separator), "", 1)

			sourceCodeRaw, err := ioutil.ReadFile(filename)
			// ファイル読み取り失敗
			if err != nil {
				println(err.Error())
				continue
			}

			language := LanguageType(filename)
			if language == "plain" || !accepted(languages, language) {
				continue
			}

			testTarget := new(TestTarget)
			testTarget.Name = name
			testTarget.Language = language
			testTarget.SourceCode = string(sourceCodeRaw)

			length++
			tmp1 = append(tmp1, testTarget)
		}
		tmp0 = append(tmp0, tmp1)
	}

	testTargets := make([]*TestTarget, 0, length)
	for _, tmp := range tmp0 {
		testTargets = append(testTargets, tmp...)
	}

	return testTargets
}

func (testTarget *TestTarget) Exec(client *client.Client, testCase *TestCase) *Detail {
	detail := new(Detail)
	detail.TestCase = testCase.Name

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := client.Do(testTarget.SourceCode, testTarget.Language, testCase.Input)

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

	detail.Time = resp.Time
	detail.Error = resp.STDERR
	detail.Output = resp.STDOUT

	return detail
}

func accepted(array []string, element string) bool {
	if len(array) == 0 {
		return true
	}

	for _, e := range array {
		if e == element {
			return true
		}
	}

	return false
}
