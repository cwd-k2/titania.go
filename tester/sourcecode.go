package tester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

// SourceCode
// contains source code, its language
type SourceCode struct {
	Name       string
	Language   string
	SourceCode string
}

// returns []*SourceCode
func MakeSourceCode(
	basepath string,
	languages []string,
	sourceCodeDirectories []string) []*SourceCode {

	tmp0 := make([][]*SourceCode, 0, len(sourceCodeDirectories))
	length := 0

	for _, dirname := range sourceCodeDirectories {
		// ソースファイル
		pattern := filepath.Join(basepath, dirname, "*.*")
		filenames, err := filepath.Glob(pattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err.Error())
			continue
		}

		tmp1 := make([]*SourceCode, 0, len(filenames))

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

			sourceCode := new(SourceCode)
			sourceCode.Name = name
			sourceCode.Language = language
			sourceCode.SourceCode = string(sourceCodeRaw)

			length++
			tmp1 = append(tmp1, sourceCode)
		}
		tmp0 = append(tmp0, tmp1)
	}

	sourceCodes := make([]*SourceCode, 0, length)
	for _, tmp := range tmp0 {
		sourceCodes = append(sourceCodes, tmp...)
	}

	return sourceCodes
}

func (sourceCode *SourceCode) Exec(client *client.Client, testCase *TestCase) *Detail {
	detail := new(Detail)
	detail.TestCase = testCase.Name

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := client.Do(sourceCode.SourceCode, sourceCode.Language, testCase.Input)

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
