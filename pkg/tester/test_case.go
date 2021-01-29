package tester

import (
	"io/ioutil"
	"path/filepath"
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
	testCases := make([]*TestCase, 0)

	for _, config := range configs {
		// 出力(正解)ファイル
		pattern := filepath.Join(basepath, config.Directory, "*"+config.OutputExtention)
		filenames, err := filepath.Glob(pattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err.Error())
			continue
		}

		for _, afile := range filenames {
			// 想定する出力があるものに対してして入力を設定する
			// 出力から先に決める
			answerBS, err := ioutil.ReadFile(afile)
			if err != nil {
				println(err.Error())
				continue
			}

			// 入力ファイル
			ifile := afile[0:len(afile)-len(config.OutputExtention)] + config.InputExtention
			inputBS, err := ioutil.ReadFile(ifile)
			if err != nil {
				println(err.Error())
				continue
			}

			name, err := filepath.Rel(basepath, afile)
			if err != nil {
				name = afile
			}
			name = name[0 : len(name)-len(config.OutputExtention)]

			testCases = append(testCases, &TestCase{name, string(inputBS), string(answerBS)})
		}

	}

	return testCases
}
