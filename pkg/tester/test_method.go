package tester

import (
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type TestMethod struct {
	Name       string
	Language   string
	OnExit     int
	FileName   string
	InputOrder []string
}

type TestMethodConfig struct {
	FileName   string   `json:"file_name"`
	OnExit     int      `json:"on_exit"`     // on_exit: 0, ...
	InputOrder []string `json:"input_order"` // ex: [stdout, stderr, input, ...]
}

// if error occurred, then nil will be returned.
func ReadTestMethod(basepath string, config TestMethodConfig) *TestMethod {
	if config.FileName == "" {
		return nil
	}

	filename := filepath.Join(basepath, config.FileName)

	name, err := filepath.Rel(basepath, filename)
	if err != nil {
		name = filename
	}

	language := runner.LangType(filename)
	if language == "plain" {
		logger.Println("Invalid test method.")
		return nil
	}

	var inputorder []string
	// TODO: validation
	if len(config.InputOrder) != 0 {
		inputorder = config.InputOrder
	} else {
		inputorder = []string{"stdout", "input", "answer"}
	}

	return &TestMethod{
		Name:       name,
		Language:   language,
		OnExit:     config.OnExit,
		FileName:   filename,
		InputOrder: inputorder,
	}
}
