package tester

import (
	"io/ioutil"
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

// returns []*TestCodes
func MakeTestCodes(
	basepath string,
	languageList []string,
	SourceCodeDirectories []string) []*TestCode {

	tmp0 := make([][]*TestCode, 0, len(SourceCodeDirectories))
	length := 0

	for _, dirname := range SourceCodeDirectories {
		// ソースファイル
		pattern := filepath.Join(basepath, dirname, "*.*")
		filenames, err := filepath.Glob(pattern)
		// ここのエラーは bad pattern
		if err != nil {
			println(err)
			continue
		}

		tmp1 := make([]*TestCode, 0, len(filenames))

		for _, filename := range filenames {
			sourceCode, err := ioutil.ReadFile(filename)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
				continue
			}

			language := LanguageType(filename)
			if language == "plain" || !accepted(languageList, language) {
				continue
			}

			name := filepath.Join(filepath.Base(basepath), strings.Replace(filename, basepath, "", 1))

			testCode := new(TestCode)
			testCode.Name = name
			testCode.Language = language
			testCode.SourceCode = string(sourceCode)

			length++
			tmp1 = append(tmp1, testCode)
		}

		tmp0 = append(tmp0, tmp1)
	}

	testCodes := make([]*TestCode, 0, length)
	for _, tmp := range tmp0 {
		testCodes = append(testCodes, tmp...)
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
