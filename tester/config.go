package tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// Config
// test configs
type Config struct {
	Host                    string   `json:"host"`
	APIKey                  string   `json:"api_key"`
	SourceCodeDirectories   []string `json:"source_code_directories"`
	TestCaseDirectories     []string `json:"test_case_directories"`
	TestCaseInputExtension  string   `json:"test_case_input_extension"`
	TestCaseOutputExtension string   `json:"test_case_output_extension"`
}

func NewConfig(baseDirectoryPath string) *Config {
	// ディレクトリ直下に titania.json がいるか確認したい
	configFileName := path.Join(baseDirectoryPath, "titania.json")
	if match, _ := filepath.Glob(configFileName); len(match) == 0 {
		return nil
	}

	// ディレクトリ直下の titania.json を読んで設定を作る
	configRawData, err := ioutil.ReadFile(configFileName)

	// File Read 失敗
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read %s.\n", configFileName)
		return nil
	}

	// ようやく設定の構造体を作れる
	config := new(Config)

	// JSON パース失敗
	if err := json.Unmarshal(configRawData, config); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse %s.\n%s\n", configFileName, err)
		return nil
	}

	return config
}
