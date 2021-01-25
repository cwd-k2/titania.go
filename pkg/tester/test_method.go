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
	SourceCode string
}

type TestMethodConfig struct {
	FileName string `json:"file_name"`
}

func NewTestMethod(basepath string, config TestMethodConfig) *TestMethod {
	if config.FileName == "" {
		return nil
	}

	filename := filepath.Join(basepath, config.FileName)

	name := strings.Replace(filename, basepath+string(filepath.Separator), "", 1)

	sourceCodeBS, err := ioutil.ReadFile(filename)
	if err != nil {
		println(err.Error())
		return nil
	}

	language := langtype.LangType(filename)
	if language == "plain" {
		println("Invalid test method.")
		return nil
	}

	return &TestMethod{name, language, string(sourceCodeBS)}
}
