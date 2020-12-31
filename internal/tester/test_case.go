package tester

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TestCase
// contains input and output texts
type TestCase struct {
	Name   string
	Input  *bytes.Buffer
	Answer *bytes.Buffer
}

type TestCaseConfig struct {
	Directory       string `json:"directory"`
	InputExtention  string `json:"input_extension"`
	OutputExtention string `json:"output_extension"`
}

// returns []*TestCases
func MakeTestCases(basepath string, configs []TestCaseConfig) []*TestCase {
	length := 0

	wg0 := &sync.WaitGroup{}
	wg0.Add(len(configs))
	tmp0 := make([][]*TestCase, len(configs))
	for i, config := range configs {
		go func(i int, config TestCaseConfig) {
			defer wg0.Done()
			// 出力(正解)ファイル
			pattern := filepath.Join(basepath, config.Directory, "*"+config.OutputExtention)
			filenames, err := filepath.Glob(pattern)
			// ここのエラーは bad pattern
			if err != nil {
				println(err.Error())
				return
			}

			wg1 := &sync.WaitGroup{}
			wg1.Add(len(filenames))
			tmp1 := make([]*TestCase, len(filenames))
			for j, answerfile := range filenames {
				go func(j int, answerfile string) {
					defer wg1.Done()
					name := mkCaseName(basepath, answerfile, config.OutputExtention)
					// 想定する出力があるものに対してして入力を設定する
					// 出力から先に決める
					answerFD, err := os.Open(answerfile)
					// ファイル読み取り失敗
					if err != nil {
						println(err.Error())
						return
					}
					defer answerFD.Close()

					// 入力ファイル
					inputfile := filepath.Join(basepath, name+config.InputExtention)

					inputFD, err := os.Open(inputfile)
					if err != nil {
						println(err.Error())
						return
					}
					defer inputFD.Close()

					input := bytes.NewBuffer(nil)
					answer := bytes.NewBuffer(nil)
					io.Copy(input, inputFD)
					io.Copy(answer, answerFD)

					length++
					tmp1[j] = &TestCase{name, input, answer}
				}(j, answerfile)
			}
			wg1.Wait()
			tmp0[i] = tmp1
		}(i, config)
	}
	wg0.Wait()

	// flatten
	testCases := make([]*TestCase, 0, length)
	for _, tmp := range tmp0 {
		for _, t := range tmp {
			if t != nil {
				testCases = append(testCases, t)
			}
		}
	}

	return testCases
}

// helper function
func mkCaseName(basepath, filename, ext string) string {
	return strings.Replace(strings.Replace(filename, basepath+string(filepath.Separator), "", 1), ext, "", 1)
}
