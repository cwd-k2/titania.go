package tester

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/cwd-k2/titania.go/pkg/pretty"
)

type TestCase struct {
	Name       string
	InputFile  string
	AnswerFile string
}

// instruction to create test cases
type TestCaseConfig struct {
	Directory          string `json:"directory"`
	InputPrefix        string `json:"input_prefix"`
	InputSuffix        string `json:"input_suffix"`
	AnswerPrefix       string `json:"answer_prefix"`
	AnswerSuffix       string `json:"answer_suffix"`
	CompatInputSuffix  string `json:"input_extension"`  // deprecated
	CompatAnswerSuffix string `json:"output_extension"` // deprecated
}

// This can return an empty slice.
// All errors are logged but ignored.
func ReadTestCases(basepath string, configs []TestCaseConfig) []*TestCase {
	tcases := make([]*TestCase, 0)

	for _, config := range configs {
		// compatibility
		if len(config.InputSuffix) == 0 && len(config.CompatInputSuffix) != 0 {
			logger.Println(pretty.Deprecated(`"input_extension"`, `"input_suffix"`))
			config.InputSuffix = config.CompatInputSuffix
		}
		if len(config.AnswerSuffix) == 0 && len(config.CompatAnswerSuffix) != 0 {
			logger.Println(pretty.Deprecated(`"output_extension"`, `"answer_suffix"`))
			config.AnswerSuffix = config.CompatAnswerSuffix
		}

		directories, err := filepath.Glob(filepath.Join(basepath, config.Directory))
		// ここのエラーは bad pattern
		if err != nil {
			logger.Printf("%+v\n", err)
			continue
		}

		for _, directory := range directories {
			inputFileNamePattern := filepath.Join(directory, config.InputPrefix+"*"+config.InputSuffix)
			inputFileNames, err := filepath.Glob(inputFileNamePattern)
			if err != nil {
				logger.Printf("%+v\n", err)
				continue
			}

			// 入力があるものに対してして出力を設定する
			for _, inputFileName := range inputFileNames {
				bname := inputFileName[len(directory)+len(config.InputPrefix)+1 : len(inputFileName)-len(config.InputSuffix)]
				// 想定する出力ファイル
				answerFileName := filepath.Join(directory, config.AnswerPrefix+bname+config.AnswerSuffix)

				name, err := filepath.Rel(basepath, filepath.Join(directory, bname))
				if err != nil {
					logger.Printf("%+v\n", err)
					continue
				}

				tcase := &TestCase{
					Name:       name,
					InputFile:  inputFileName,
					AnswerFile: answerFileName,
				}

				tcases = append(tcases, tcase)
			}
		}
	}

	return tcases
}

func (t *TestCase) WriteInputDataTo(w io.Writer) error {
	fp, err := os.Open(t.InputFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}

func (t *TestCase) WriteAnswerDataTo(w io.Writer) error {
	fp, err := os.Open(t.AnswerFile)
	if err != nil {
		return err
	}
	if _, err := bufio.NewReader(fp).WriteTo(w); err != nil {
		return err
	}
	return fp.Close()
}
