package tester

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cwd-k2/titania.go/internal/pkg/langtype"
)

// TestTarget
// contains source code, its language
type TestTarget struct {
	Name       string
	Language   string
	SourceCode string
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
	wg0.Add(len(configs))

	tmp0 := make([][]*TestTarget, len(configs))

	for i, config := range configs {
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
			wg1.Add(len(filenames))

			tmp1 := make([]*TestTarget, len(filenames))

			for j, filename := range filenames {
				go func(j int, filename string) {
					defer wg1.Done()
					name := strings.Replace(filename, basepath+string(filepath.Separator), "", 1)

					sourceCodeBS, err := ioutil.ReadFile(filename)
					// ファイル読み取り失敗
					if err != nil {
						println(err.Error())
						return
					}

					language := langtype.LangType(filename)
					if language == "plain" || !acceptable(languages, language) {
						return
					}

					length++
					tmp1[j] = &TestTarget{name, language, string(sourceCodeBS), expect}
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
		testTargets = append(testTargets, tmp...)
	}

	return testTargets
}

func acceptable(array []string, element string) bool {
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
