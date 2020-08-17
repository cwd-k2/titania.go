package tester

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// TestCode
// contains source code, its language
type TestCode struct {
	Name       string
	Language   string
	SourceCode string
}

// returns map[string]*TestCodes
func MakeTestCodes(
	baseDirectoryPath string,
	languageList []string,
	SourceCodeDirectories []string) map[string]*TestCode {

	testCodes := make(map[string]*TestCode)

	for _, dirname := range SourceCodeDirectories {
		// ソースファイル
		sourceFileNamePattern := path.Join(baseDirectoryPath, dirname, "*.*")
		sourceFileNames, err := filepath.Glob(sourceFileNamePattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err)
			continue
		}

		for _, sourceFileName := range sourceFileNames {
			byteArray, err := ioutil.ReadFile(sourceFileName)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
				continue
			}

			language := LanguageType(sourceFileName)
			if language == "plain" || !accepted(languageList, language) {
				continue
			}

			name := path.Join(
				filepath.Base(baseDirectoryPath),
				strings.Replace(sourceFileName, baseDirectoryPath, "", 1))

			testCode := new(TestCode)
			testCode.Name = name
			testCode.SourceCode = string(byteArray)
			testCode.Language = language

			testCodes[name] = testCode

		}
	}

	return testCodes
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
