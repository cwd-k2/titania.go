package tester

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/cwd-k2/titania.go/internal/client"
	"github.com/cwd-k2/titania.go/pkg/pretty"
)

// Config
// test configs
type Config struct {
	ClientConfig client.Config      `json:"client"`
	TestTarget   []TestTargetConfig `json:"test_target"`
	TestCase     []TestCaseConfig   `json:"test_case"`
	TestMethod   TestMethodConfig   `json:"test_method"`
}

func NewConfig(basepath string) *Config {
	// ディレクトリ直下に titania.json がいるか確認したい
	filename := path.Join(basepath, "titania.json")
	if match, _ := filepath.Glob(filename); len(match) == 0 {
		return nil
	}

	// ディレクトリ直下の titania.json を読んで設定を作る
	rawData, err := ioutil.ReadFile(filename)

	// File Read 失敗
	if err != nil {
		pretty.Printf("Couldn't read %s.\n", filename)
		return nil
	}

	// ようやく設定の構造体を作れる
	config := new(Config)

	// JSON パース失敗
	if err := json.Unmarshal(rawData, config); err != nil {
		pretty.Printf("Couldn't parse %s.\n%s\n", filename, err)
		return nil
	}

	return config
}
