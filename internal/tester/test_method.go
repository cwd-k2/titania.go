package tester

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type TestMethod struct {
	Name       string
	Language   string
	SourceCode *bytes.Buffer
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

	sourceCodeFD, err := os.Open(filename)
	if err != nil {
		println(err.Error())
		return nil
	}
	defer sourceCodeFD.Close()
	sourceCode := bytes.NewBuffer(nil)
	io.Copy(sourceCode, sourceCodeFD)

	language := LanguageType(filename)
	if language == "plain" {
		println("Invalid test method.")
		return nil
	}

	return &TestMethod{name, language, sourceCode}
}
