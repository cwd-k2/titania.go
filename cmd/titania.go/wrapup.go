package main

import (
	"bufio"
	"encoding/json"
	"os"

	. "github.com/cwd-k2/titania.go/pkg/pretty"
	"github.com/cwd-k2/titania.go/pkg/simplejson"
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func final(uresults []*tester.TestUnitResult) {
	Printf("\n%s\n", Bold("ALL DONE"))

	for _, uresult := range uresults {
		showTestUnitResult(uresult)
	}
}

func showTestUnitResult(uresult *tester.TestUnitResult) {
	Printf("\n%s\n", Bold(Cyan(uresult.Name)))

	for _, tresult := range uresult.TestTargets {
		showTestTargetResult(tresult)
	}
}

func showTestTargetResult(tresult *tester.TestTargetResult) {
	Printf("%s: %s\n", Bold(tresult.Language), Bold(Blue(tresult.Name)))

	for _, cresult := range tresult.TestCases {
		showTestCaseResult(cresult)
	}
}

func showTestCaseResult(cresult *tester.TestCaseResult) {
	switch cresult.Result {
	case "CLIENT ERROR":
		Printf("%s: %s\n", Magenta(cresult.Name), Magenta(cresult.Result))
	case "SERVER ERROR":
		Printf("%s: %s\n", Blue(cresult.Name), Blue(cresult.Result))
	case "TESTER ERROR":
		Printf("%s: %s\n", Bold(Red(cresult.Name)), Bold(Red(cresult.Result)))
	case "PASS":
		fallthrough
	case "FAIL":
		if cresult.IsExpected {
			Printf("%s: %s %ss\n", Green(cresult.Name), Green(cresult.Result), cresult.Time)
		} else {
			Printf("%s: %s %ss\n", Yellow(cresult.Name), Yellow(cresult.Result), cresult.Time)
		}
	default:
		if cresult.IsExpected {
			Printf("%s: %s\n", Green(cresult.Name), Green(cresult.Result))
		} else {
			Printf("%s: %s\n", Red(cresult.Name), Red(cresult.Result))
		}
	}
}

func printjson(uresults []*tester.TestUnitResult) {
	builder := simplejson.NewArrayBuilder()
	for _, uresult := range uresults {
		builder.AddObject(func(obj simplejson.Object) {
			obj.SetString("name", uresult.Name)
			obj.SetString("test_method", uresult.TestMethod)
			obj.SetArray("fruits", func(arr simplejson.Array) {
				for _, tresult := range uresult.TestTargets {
					arr.AddObject(func(obj simplejson.Object) {
						obj.SetString("name", tresult.Name)
						obj.SetString("language", tresult.Language)
						obj.SetString("expect", tresult.Expect)
						obj.SetArray("details", func(arr simplejson.Array) {
							for _, cresult := range tresult.TestCases {
								arr.AddObject(func(obj simplejson.Object) {
									obj.SetString("test_case", cresult.Name)
									obj.SetString("result", cresult.Result)
									obj.SetBool("is_expected", cresult.IsExpected)
									obj.SetString("time", cresult.Time)
									obj.SetStringFromFile("output", cresult.Output)
									obj.SetStringFromFiles("others", cresult.Others, "")
									obj.SetString("error", cresult.Error)
								})
							}
						})
					})
				}
			})
		})
	}
	stdoutbuf := bufio.NewWriter(os.Stdout)
	defer stdoutbuf.Flush()

	if prettyprint {
		o := make([]*struct {
			Name       string `json:"name"`
			TestMethod string `json:"test_method"`
			Fruits     []*struct {
				TestTarget string `json:"test_target"`
				Language   string `json:"language"`
				Expect     string `json:"expect"`
				Details    []*struct {
					TestCase   string `json:"test_case"`
					Result     string `json:"result"`
					IsExpected bool   `json:"is_expected"`
					Time       string `json:"time"`
					Output     string `json:"output"`
					Others     string `json:"others"`
					Error      string `json:"error"`
				} `json:"details"`
			} `json:"fruits"`
		}, 0)

		dec := json.NewDecoder(builder.Build())
		dec.Decode(&o)

		enc := json.NewEncoder(stdoutbuf)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)

		if err := enc.Encode(o); err != nil {
			panic(err)
		}

	} else {
		stdoutbuf.ReadFrom(builder.Build())
	}

}
