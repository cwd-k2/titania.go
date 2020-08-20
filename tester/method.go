package tester

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/cwd-k2/titania.go/client"
)

type Method struct {
	Name       string
	Language   string
	SourceCode string
}

func NewMethod(basepath, testMethodFileName string) *Method {
	if testMethodFileName == "" {
		return nil
	}

	filename := filepath.Join(basepath, testMethodFileName)

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

	method := new(Method)
	method.Name = name
	method.Language = language
	method.SourceCode = string(sourceCodeRaw)

	return method
}

func (method *Method) Exec(client *client.Client, testCase *TestCase, detail *Detail) (string, string) {

	// input for test_method goes in this format.
	elems := []string{
		testCase.Input,
		"END_OF_CHUNK",
		detail.Output,
		"END_OF_CHUNK",
		testCase.Answer,
		"END_OF_CHUNK",
	}

	input := strings.Join(elems, "")

	// 実際に paiza.io の API を利用して実行結果をもらう
	resp, err := client.Do(method.SourceCode, method.Language, input)

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
