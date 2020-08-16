package tester

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// TestUnit
// contains source code, its language
type TestUnit struct {
	SourceCode string
	Language   string
}

// returns map[string]*TestUnits
func MakeTestUnits(
	baseDirectoryPath string,
	languageList []string,
	sourceCodeDirectories []string) map[string]*TestUnit {

	testUnits := make(map[string]*TestUnit)

	for _, dirname := range sourceCodeDirectories {
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

			unitName :=
				path.Join(
					filepath.Base(baseDirectoryPath),
					strings.Replace(sourceFileName, baseDirectoryPath, "", 1))

			testUnits[unitName] = new(TestUnit)
			testUnits[unitName].SourceCode = string(byteArray)
			testUnits[unitName].Language = language

		}
	}

	return testUnits
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
