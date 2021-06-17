package tester

import (
	"os"
	"path/filepath"
)

type TestCase struct {
	Name       string
	InputData  []byte
	AnswerData []byte
}

type TestCaseConfig struct {
	Directory       string `json:"directory"`
	InputExtention  string `json:"input_extension"`
	OutputExtention string `json:"output_extension"`
}

// This can return an empty slice.
// All errors are logged but ignored.
func ReadTestCases(basepath string, configs []TestCaseConfig) []*TestCase {
	tcases := make([]*TestCase, 0)

	for _, config := range configs {
		// 入力ファイル
		pattern := filepath.Join(basepath, config.Directory, "*"+config.InputExtention)
		inputFileNames, err := filepath.Glob(pattern)
		// ここのエラーは bad pattern
		if err != nil {
			logger.Printf("%+v\n", err)
			continue
		}

		for _, inputFileName := range inputFileNames {
			// 入力があるものに対してして出力を設定する
			// 入力から先に決める

			// 想定する出力ファイル
			answerFileName := inputFileName[0:len(inputFileName)-len(config.InputExtention)] + config.OutputExtention

			name, _ := filepath.Rel(basepath, inputFileName)
			name = name[0 : len(name)-len(config.InputExtention)]

			tcase := &TestCase{
				Name: name,
			}

			if tcase.InputData, err = os.ReadFile(inputFileName); err != nil {
				logger.Printf("%+v\n", err)
				continue
			}
			// 出力ファイルはなくてもよい
			if tcase.AnswerData, err = os.ReadFile(answerFileName); err != nil {
				logger.Printf("%+v\n", err)
			}

			tcases = append(tcases, tcase)
		}

	}

	return tcases
}
