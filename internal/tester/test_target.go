package tester

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TestTarget
// contains source code, its language
type TestTarget struct {
	Name       string
	Language   string
	SourceCode *bytes.Buffer
	Expect     string
}

type TestTargetConfig struct {
	Pattern string `json:"pattern"`
	Expect  string `json:"expect"`
}

// returns []*TestTarget
func MakeTestTargets(basepath string, languages []string, configs []TestTargetConfig) []*TestTarget {
	length := 0

	wg0 := &sync.WaitGroup{}
	tmp0 := make([][]*TestTarget, len(configs))
	for i, config := range configs {
		wg0.Add(1)
		// ソースファイル
		go func(i int, config TestTargetConfig) {
			defer wg0.Done()

			pattern := filepath.Join(basepath, config.Pattern)
			filenames, err := filepath.Glob(pattern)
			// ここのエラーは bad pattern
			if err != nil {
				println(err.Error())
				return
			}

			var expect string
			if config.Expect != "" {
				expect = config.Expect
			} else {
				expect = "PASS"
			}

			wg1 := &sync.WaitGroup{}
			tmp1 := make([]*TestTarget, len(filenames))
			for j, filename := range filenames {
				wg1.Add(1)

				go func(j int, filename string) {
					defer wg1.Done()
					name := strings.Replace(filename, basepath+string(filepath.Separator), "", 1)

					sourceCodeFD, err := os.Open(filename)
					// ファイル読み取り失敗
					if err != nil {
						println(err.Error())
						return
					}
					defer sourceCodeFD.Close()
					sourceCode := bytes.NewBuffer(nil)
					io.Copy(sourceCode, sourceCodeFD)

					language := LanguageType(filename)
					if language == "plain" || !accepted(languages, language) {
						return
					}

					length++
					tmp1[j] = &TestTarget{name, language, sourceCode, expect}
				}(j, filename)
			}
			wg1.Wait()
			tmp0[i] = tmp1
		}(i, config)
	}
	wg0.Wait()

	// flatten
	testTargets := make([]*TestTarget, 0, length)
	for _, tmp := range tmp0 {
		for _, t := range tmp {
			if t != nil {
				testTargets = append(testTargets, t)
			}
		}
	}

	return testTargets
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
