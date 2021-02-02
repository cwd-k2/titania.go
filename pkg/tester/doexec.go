package tester

import (
	"fmt"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/paizaio"
)

type singleresult struct {
	Result      string
	Time        string
	BuildSTDOUT string
	BuildSTDERR string
	STDOUT      string
	STDERR      string
	Error       string
}

func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *Detail {
	// TODO: refactoring
	// fire paiza.io API
	sres1 := t.do(target.Language, target.SourceCode, tcase.Input)

	result := sres1.Result
	// anything but BuildSTDOUT or STDOUT
	// TODO: still ignoring BuildSTDOUT...
	// TODO: Detail に Build 時のものとか他に何も入れてないからちょっと辛い
	errstr := sres1.BuildSTDERR + sres1.STDERR + sres1.Error

	// making result string
	// if not confirmed yet
	if len(result) == 0 {
		// TODO: Method Execution `on` specified result.
		if t.TestMethod != nil {
			// input for test_method goes in this format.
			// output + "\0" + input + "\0" + answer
			// TODO: the order and element should be specified by config.
			input := strings.Join([]string{sres1.STDOUT, *tcase.Input, *tcase.Answer}, "\000")

			// TestMethod
			sres2 := t.do(t.TestMethod.Language, t.TestMethod.SourceCode, &input)

			// if not confirmed yet
			if len(sres2.Result) == 0 {
				// mainly expecting PASS or FAIL
				result = strings.TrimRight(sres2.STDOUT, "\n")
			} else {
				result = fmt.Sprintf("METHOD %s", sres2.Result)
			}

			// still anything but STDOUT
			errstr += sres2.BuildSTDERR + sres2.STDERR + sres2.Error
		} else {
			// simple comparison
			if sres1.STDOUT == *tcase.Answer {
				result = "PASS"
			} else {
				result = "FAIL"
			}

		}
	}

	return &Detail{tcase.Name, result, result == target.Expect, sres1.Time, sres1.STDOUT, errstr}
}

// errors are treated as string
func (t *TestUnit) do(language string, sourceCode, input *string) *singleresult {
	// TODO: how can I treat build time?

	req1 := &paizaio.RunnersCreateRequest{
		Language:        language,
		SourceCode:      *sourceCode,
		Input:           *input,
		Longpoll:        true,
		LongpollTimeout: 16,
	}

	res1, err := t.Client.RunnersCreate(req1)
	if err != nil {
		handle(err)
	}

	req2 := &paizaio.RunnersGetDetailsRequest{
		ID: res1.ID,
	}

	res2, err := t.Client.RunnersGetDetails(req2)
	if err != nil {
		return handle(err)
	}

	if res2.BuildExitCode != 0 {
		result := fmt.Sprintf("BUILD %s", strings.ToUpper(res2.BuildResult))
		return &singleresult{
			Result:      result,
			Time:        res2.BuildTime,
			BuildSTDOUT: res2.BuildSTDOUT,
			BuildSTDERR: res2.BuildSTDERR,
		}
	}

	if res2.ExitCode != 0 || res2.Result != "success" {
		result := fmt.Sprintf("EXECUTION %s", strings.ToUpper(res2.Result))
		return &singleresult{
			Result:      result,
			Time:        res2.Time,
			BuildSTDOUT: res2.BuildSTDOUT,
			BuildSTDERR: res2.BuildSTDERR,
			STDOUT:      res2.STDOUT,
			STDERR:      res2.STDERR,
		}
	}

	return &singleresult{
		Time:        res2.Time,
		BuildSTDOUT: res2.BuildSTDOUT,
		BuildSTDERR: res2.BuildSTDERR,
		STDOUT:      res2.STDOUT,
		STDERR:      res2.STDERR,
	}
}

func handle(err error) *singleresult {
	var result, errstr string

	switch err := err.(type) {
	case paizaio.ServerError:
		result = "SERVER ERROR"
		errstr = fmt.Sprintf("HTTP response status code: %d\n%s", err.Code, err.Error())
	case paizaio.ClientError:
		result = "CLIENT ERROR"
		errstr = fmt.Sprintf("HTTP response status code: %d\n%s", err.Code, err.Error())
	case paizaio.RunnerError:
		result = "RUNNER ERROR"
		errstr = fmt.Sprintf("Error occurred at paiza.io code runner.\n%s", err.Error())
	default:
		result = "TESTER ERROR"
		errstr = err.Error()
	}

	return &singleresult{Result: result, Error: errstr}
}
