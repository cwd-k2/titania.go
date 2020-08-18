package tester

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// TestCase
// contains input and output texts
type TestCase struct {
	Name   string
	Input  string
	Output string
}

// returns map[string]*TestCases
func MakeTestCases(
	basepath string,
	testCaseDirectories []string,
	inputExt, outputExt string) map[string]*TestCase {

	testCases := make(map[string]*TestCase)

	for _, dirname := range testCaseDirectories {
		// 出力(正解)ファイル
		outFileNamePattern := path.Join(basepath, dirname, "*"+outputExt)
		outFileNames, err := filepath.Glob(outFileNamePattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err)
			continue
		}

		// 入力ファイル
		inFileNamePattern := path.Join(basepath, dirname, "*"+inputExt)
		inFileNames, err := filepath.Glob(inFileNamePattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err)
			continue
		}

		// 出力から先に決める
		for _, filename := range outFileNames {
			name := makeCaseName(basepath, filename, outputExt)

			byteArray, err := ioutil.ReadFile(filename)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
				testCases[name] = nil
				continue
			}

			testCase := new(TestCase)
			testCase.Name = name
			testCase.Output = string(byteArray)

			testCases[name] = testCase
		}

		// 想定する出力があるものに大して入力を設定する
		for _, filename := range inFileNames {
			name := makeCaseName(basepath, filename, inputExt)
			testCase := testCases[name]
			// 出力が用意されてなかったら作りません
			if testCase == nil {
				continue
			}

			byteArray, err := ioutil.ReadFile(filename)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
				testCase = nil
				continue
			}

			testCase.Input = string(byteArray)
		}
	}

	return testCases

}

// helper function
func makeCaseName(basepath, filename, ext string) string {
	return path.Join(
		filepath.Base(basepath),
		strings.Replace(strings.Replace(filename, basepath, "", 1), ext, "", 1))
}
