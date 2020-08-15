package tester

import (
	"fmt"
	"io/ioutil"
	"os"
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
	sourceCodeDirectories []string) map[string]*TestUnit {

	testUnits := make(map[string]*TestUnit)

	for _, dirname := range sourceCodeDirectories {
		// ソースファイル
		sourceFileNamePattern := path.Join(baseDirectoryPath, dirname, "*.*")
		sourceFileNames, err := filepath.Glob(sourceFileNamePattern)
		// ここのエラーは bad pattern
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		for _, sourceFileName := range sourceFileNames {
			byteArray, err := ioutil.ReadFile(sourceFileName)
			// ファイル読み取り失敗
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}

			language := LanguageType(sourceFileName)
			if language == "plain" {
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
