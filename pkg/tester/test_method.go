package tester

import (
	"io/ioutil"
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/runner"
)

type TestMethod struct {
	Name       string
	Language   string
	CodeData   []byte
	Delimiter  string
	OnExit     int
	InputOrder []string
}

type TestMethodConfig struct {
	FileName   string   `json:"file_name"`
	OnExit     int      `json:"on_exit"`     // on_exit: 0, ...
	Delimiter  string   `json:"delimiter"`   // ex: 9E806A3A-8E0C-4CF4-8139-4ABCC2443E4E
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

	var delimiter string
	if len(config.Delimiter) != 0 {
		delimiter = config.Delimiter
	} else {
		delimiter = "\x00"
	}

	var inputorder []string
	// TODO: validation
	if len(config.InputOrder) != 0 {
		inputorder = config.InputOrder
	} else {
		inputorder = []string{"stdout", "newline", "input", "newline", "answer"}
	}

	tmethod := &TestMethod{
		Name:       name,
		Language:   language,
		Delimiter:  delimiter,
		OnExit:     config.OnExit,
		InputOrder: inputorder,
	}

	if tmethod.CodeData, err = ioutil.ReadFile(filename); err != nil {
		logger.Printf("%+v\n", err)
		return nil
	}

	return tmethod
}
