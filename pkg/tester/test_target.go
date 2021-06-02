package tester

import (
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type TestTarget struct {
	Name     string
	Language string
	FileName string
	Expect   string
}

type TestTargetConfig struct {
	Pattern string `json:"pattern"`
	Expect  string `json:"expect"`
}

// This can return an empty slice.
// All errors are logged but ignored.
func ReadTestTargets(basepath string, configs []TestTargetConfig) []*TestTarget {
	targets := make([]*TestTarget, 0)

	for _, config := range configs {
		filenames, err := filepath.Glob(filepath.Join(basepath, config.Pattern))
		if err != nil {
			logger.Printf("%+v\n", err)
			continue
		}

		expect := config.Expect
		if len(expect) == 0 {
			expect = "PASS"
		}

		for _, filename := range filenames {
			language := runner.LangType(filename)
			if language == "plain" || len(languages) > 0 && !acceptable(languages, language) {
				continue
			}

			name, err := filepath.Rel(basepath, filename)
			if err != nil {
				name = filename
			}

			targets = append(targets, &TestTarget{name, language, filename, expect})
		}
	}

	return targets
}

func acceptable(languages []string, language string) bool {
	for _, e := range languages {
		if e == language {
			return true
		}
	}

	return false
}
