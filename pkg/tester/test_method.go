package tester

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/internal/pkg/langtype"
)

type TestMethod struct {
	Name       string
	Language   string
	OnResult   string
	SourceCode []byte
	InputOrder []string
}

type TestMethodConfig struct {
	FileName   string   `json:"file_name"`
	On         string   `json:"on"`          // one of SUCCESS, EXECUTION FAILURE, ...
	InputOrder []string `json:"input_order"` // ex: [stdout, stderr, input, ...]
}

// Creates TestMethod struct.
// if error occurred, then nil will be returned.
// TODO: error handling.
func NewTestMethod(basepath string, config TestMethodConfig) *TestMethod {
	if config.FileName == "" {
		return nil
	}

	filename := filepath.Join(basepath, config.FileName)

	name, err := filepath.Rel(basepath, filename)
	if err != nil {
		name = filename
	}

	sourceCodeBS, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Printf("%+v\n", err)
		return nil
	}

	language := langtype.LangType(filename)
	if language == "plain" {
		logger.Println("Invalid test method.")
		return nil
	}

	var onresult string
	if len(config.On) != 0 {
		onresult = strings.ToUpper(config.On)
	} else {
		onresult = "SUCCESS"
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
		OnResult:   onresult,
		SourceCode: sourceCodeBS,
		InputOrder: inputorder,
	}
}
