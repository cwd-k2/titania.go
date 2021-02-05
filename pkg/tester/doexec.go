package tester

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/paizaio"
)

type singleresult struct {
	Result      string
	Time        string
	BuildSTDOUT []byte
	BuildSTDERR []byte
	STDOUT      []byte
	STDERR      []byte
	Error       string
}

func (t *TestUnit) exec(target *TestTarget, tcase *TestCase) *Detail {
	// TODO: refactoring
	// fire paiza.io API
	sres1 := t.do(target.Language, target.SourceCode, tcase.Input)

	var (
		result string
		// anything but BuildSTDOUT or STDOUT
		// TODO: still ignoring BuildSTDOUT...
		errstr = string(sres1.BuildSTDERR) + string(sres1.STDERR) + sres1.Error
	)

	// making result string
	// TODO: Method Execution `on` specified result.
	if t.TestMethod != nil && sres1.Result == "SUCCESS" {
		// input for test_method goes in this format.
		// output + "\0" + input + "\0" + answer
		// TODO: the order and element should be specified by config.
		var input []byte
		input = append(input, sres1.STDOUT...)
		input = append(input, '\000')
		input = append(input, tcase.Input...)
		input = append(input, '\000')
		input = append(input, tcase.Answer...)

		// TestMethod
		sres2 := t.do(t.TestMethod.Language, t.TestMethod.SourceCode, input)

		if sres2.Result == "SUCCESS" {
			result = strings.TrimRight(string(sres2.STDOUT), "\n") // mainly expecting PASS or FAIL
		} else {
			result = fmt.Sprintf("METHOD %s", sres2.Result)
		}

		errstr += string(sres2.BuildSTDERR) + string(sres2.STDERR) + sres2.Error

	} else if sres1.Result == "SUCCESS" {
		// simple comparison
		if bytes.Equal(sres1.STDOUT, tcase.Answer) {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	} else {
		result = sres1.Result
	}

	return &Detail{
		TestCase:   tcase.Name,
		Result:     result,
		IsExpected: result == target.Expect,
		Time:       sres1.Time,
		Output:     string(sres1.STDOUT),
		Error:      errstr,
	}
}

// errors are treated as string
func (t *TestUnit) do(language string, sourceCode, input []byte) *singleresult {
	// TODO: how can I treat build time?

	req1 := &paizaio.RunnersCreateRequest{
		Language:        language,
		SourceCode:      string(sourceCode),
		Input:           string(input),
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
			BuildSTDOUT: []byte(res2.BuildSTDOUT),
			BuildSTDERR: []byte(res2.BuildSTDERR),
		}
	}

	if res2.ExitCode != 0 || res2.Result != "success" {
		result := fmt.Sprintf("EXECUTION %s", strings.ToUpper(res2.Result))
		return &singleresult{
			Result:      result,
			Time:        res2.Time,
			BuildSTDOUT: []byte(res2.BuildSTDOUT),
			BuildSTDERR: []byte(res2.BuildSTDERR),
			STDOUT:      []byte(res2.STDOUT),
			STDERR:      []byte(res2.STDERR),
		}
	}

	return &singleresult{
		Result:      strings.ToUpper(res2.Result),
		Time:        res2.Time,
		BuildSTDOUT: []byte(res2.BuildSTDOUT),
		BuildSTDERR: []byte(res2.BuildSTDERR),
		STDOUT:      []byte(res2.STDOUT),
		STDERR:      []byte(res2.STDERR),
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
