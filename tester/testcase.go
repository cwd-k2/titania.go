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
	baseDirectoryPath string,
	testCaseDirectories []string,
	inputExt, outputExt string) map[string]*TestCase {

	testCases := make(map[string]*TestCase)

	for _, dirname := range testCaseDirectories {
		// 出力(正解)ファイル
		outFileNamePattern := path.Join(baseDirectoryPath, dirname, "*"+outputExt)
		outFileNames, err := filepath.Glob(outFileNamePattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err)
			continue
		}

		// 入力ファイル
		inFileNamePattern := path.Join(baseDirectoryPath, dirname, "*"+inputExt)
		inFileNames, err := filepath.Glob(inFileNamePattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err)
			continue
		}

		// 出力から先に決める
		for _, outFileName := range outFileNames {
			caseName := makeCaseName(baseDirectoryPath, outFileName, outputExt)

			byteArray, err := ioutil.ReadFile(outFileName)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
				testCases[caseName] = nil
				continue
			}

			testCases[caseName] = new(TestCase)
			testCases[caseName].Name = caseName
			testCases[caseName].Output = string(byteArray)
		}

		// 想定する出力があるものに大して入力を設定する
		for _, inFileName := range inFileNames {
			caseName := makeCaseName(baseDirectoryPath, inFileName, inputExt)

			byteArray, err := ioutil.ReadFile(inFileName)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
				testCases[caseName] = nil
				continue
			}

			// 出力が用意されてなかったら作りません
			if testCases[caseName] != nil {
				testCases[caseName].Input = string(byteArray)
			}
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
