package tester

import (
	"io/ioutil"
	"path/filepath"

	"github.com/cwd-k2/titania.go/internal/pkg/langtype"
)

type TestMethod struct {
	Name       string
	Language   string
	SourceCode []byte
}

type TestMethodConfig struct {
	FileName string `json:"file_name"`
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

	return &TestMethod{name, language, sourceCodeBS}
}
