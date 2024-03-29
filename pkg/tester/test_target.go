package tester

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type TestTarget struct {
	Name       string
	Language   string
	SourceFile string
	Expect     map[string]string
}

type TestTargetConfig struct {
	Pattern string         `json:"pattern"`
	Expect  ExpectedResult `json:"expect"`
}

type ExpectedResult struct {
	Map map[string]string
}

func (e *ExpectedResult) UnmarshalJSON(data []byte) error {
	e.Map = map[string]string{}
	var dyn interface{}
	if err := json.Unmarshal(data, &dyn); err != nil {
		return err
	}
	switch dyn.(type) {
	case string:
		e.Map["default"] = dyn.(string)
	case map[string]string:
		for k, v := range dyn.(map[string]string) {
			e.Map[k] = v
		}
	}
	return nil
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

		expect := map[string]string{"default": "PASS"}
		for k, v := range config.Expect.Map {
			expect[k] = v
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

			target := &TestTarget{
				Name:       name,
				Language:   language,
				SourceFile: filename,
				Expect:     expect,
			}

			targets = append(targets, target)
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

func (t *TestTarget) WriteSouceCodeTo(w io.Writer) error {
	fp, err := os.Open(t.SourceFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}
