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

// returns []*TestCases
func MakeTestCases(
	basepath string,
	testCaseDirectories []string,
	inputExt, answerExt string) []*TestCase {

	tmp0 := make([][]*TestCase, 0, len(testCaseDirectories))
	length := 0

	for _, dirname := range testCaseDirectories {
		// 出力(正解)ファイル
		pattern := filepath.Join(basepath, dirname, "*"+answerExt)
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
			name := mkCaseName(basepath, answerfile, answerExt)

			answer, err := ioutil.ReadFile(answerfile)
			// ファイル読み取り失敗
			if err != nil {
				println(err.Error())
				continue
			}

			// 入力ファイル
			inputfile := filepath.Join(basepath, dirname, filepath.Base(name)+inputExt)

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
