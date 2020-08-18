package tester

import (
	"io/ioutil"
	"path/filepath"
	"strings"
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
	SourceCodeDirectories []string) []*SourceCode {

	tmp0 := make([][]*SourceCode, 0, len(SourceCodeDirectories))
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

		tmp1 := make([]*SourceCode, 0, len(filenames))

		for _, filename := range filenames {
			name := filepath.Join(filepath.Base(basepath), strings.Replace(filename, basepath, "", 1))

			sourceCodeRaw, err := ioutil.ReadFile(filename)
			// ファイル読み取り失敗
			if err != nil {
				println(err)
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
