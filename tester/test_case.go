package tester

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

// TestCase
// contains input and output texts
type TestCase struct {
	Name   string
	Input  string
	Answer string
}

type TestCaseConfig struct {
	Directory       string `json:"directory"`
	InputExtention  string `json:"input_extension"`
	OutputExtention string `json:"output_extension"`
}

// returns []*TestCases
func MakeTestCases(basepath string, configs []TestCaseConfig) []*TestCase {

	tmp0 := make([][]*TestCase, 0, len(configs))
	length := 0

	for _, config := range configs {
		// 出力(正解)ファイル
		pattern := filepath.Join(basepath, config.Directory, "*"+config.OutputExtention)
		filenames, err := filepath.Glob(pattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err.Error())
			continue
		}
		tmp1 := make([]*TestCase, 0, len(filenames))

		// 想定する出力があるものに対してして入力を設定する
		// 出力から先に決める
		for _, answerfile := range filenames {
			name := mkCaseName(basepath, answerfile, config.OutputExtention)

			answer, err := ioutil.ReadFile(answerfile)
			// ファイル読み取り失敗
			if err != nil {
				println(err.Error())
				continue
			}

			// 入力ファイル
			inputfile := filepath.Join(basepath, name+config.InputExtention)

			input, err := ioutil.ReadFile(inputfile)
			if err != nil {
				println(err.Error())
				continue
			}

			testCase := new(TestCase)
			testCase.Name = name
			testCase.Input = string(input)
			testCase.Answer = string(answer)

			length++
			tmp1 = append(tmp1, testCase)
		}
		tmp0 = append(tmp0, tmp1)
	}

	// flatten
	testCases := make([]*TestCase, 0, length)
	for _, tmp := range tmp0 {
		testCases = append(testCases, tmp...)
	}

	return testCases

}

// helper function
func mkCaseName(basepath, filename, ext string) string {
	return strings.Replace(strings.Replace(filename, basepath+string(filepath.Separator), "", 1), ext, "", 1)
}
