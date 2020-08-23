package tester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
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

	sourceCodeRaw, err := ioutil.ReadFile(filename)
	if err != nil {
		println(err.Error())
		return nil
	}

	language := LanguageType(filename)
	if language == "plain" {
		println("Invalid test method.")
		return nil
	}

	name := strings.Replace(filename, basepath+string(filepath.Separator), "", 1)

	testMethod := new(TestMethod)
	testMethod.Name = name
	testMethod.Language = language
	testMethod.SourceCode = string(sourceCodeRaw)

	return testMethod
}

// Exec
// returns STDOUT and STDERR for test method execution.
// STDOUT should be the result, and STDERR should be the reason for that output.
func (testMethod *TestMethod) Exec(client *client.Client, testCase *TestCase, detail *Detail) (string, string) {

	// input for test_method goes in this format.
	// output + "\0" + input + "\0" + answer + "\0"
	elems := []string{
		detail.Output,
		"\000",
		testCase.Input,
		"\000",
		testCase.Answer,
		"\000",
	}

	input := strings.Join(elems, "")

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := client.Do(testMethod.SourceCode, testMethod.Language, input)

	// Errors that are not related to source_code
	if err != nil {
		if err.Code >= 500 {
			return "SERVER ERROR", err.Error()
		} else if err.Code >= 400 {
			return "CLIENT ERROR", err.Error()
		} else {
			return "TESTER ERROR", err.Error()
		}
	}

	// ビルドエラー
	if !(resp.BuildResult == "success" || resp.BuildResult == "") {
		return fmt.Sprintf("METHOD BUILD %s", strings.ToUpper(resp.BuildResult)), resp.BuildSTDERR
	}

	// 実行時エラー
	if resp.Result != "success" {
		return fmt.Sprintf("METHOD EXECUTION %s", strings.ToUpper(resp.Result)), resp.STDERR
	}

	// expect: "PASS" or "FAIL"
	return resp.STDOUT, resp.STDERR
}
