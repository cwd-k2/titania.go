package tester

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type Config struct {
	ClientConfig runner.Config      `json:"client"`
	TestTarget   []TestTargetConfig `json:"test_target"`
	TestCase     []TestCaseConfig   `json:"test_case"`
	TestMethod   TestMethodConfig   `json:"test_method"`
}

// if error occurred, then nil will be returned.
func ReadConfig(dirname string) *Config {
	basepath, err := filepath.Abs(dirname)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil
	}

	// ディレクトリ直下に titania.json がいるか確認したい
	filename := filepath.Join(basepath, "titania.json")
	if match, _ := filepath.Glob(filename); len(match) == 0 {
		return nil
	}

	// ディレクトリ直下の titania.json を読んで設定を作る
	rawData, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Printf("Couldn't read %s.\n%+v\n", filename, err)
		return nil
	}

	// ようやく設定の構造体を作れる
	config := &Config{}
	if err := json.Unmarshal(rawData, config); err != nil {
		logger.Printf("Couldn't parse %s.\n%+v\n", filename, err)
		return nil
	}

	return config
}
