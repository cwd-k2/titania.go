package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
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
		if cresult.Result == cresult.Expect {
			Printf("%s: %s %ss\n", Green(cresult.Name), Green(cresult.Result), cresult.Time)
		} else {
			Printf("%s: %s %ss\n", Yellow(cresult.Name), Yellow(cresult.Result), cresult.Time)
		}
	default:
		if cresult.Result == cresult.Expect {
			Printf("%s: %s\n", Green(cresult.Name), Green(cresult.Result))
		} else {
			Printf("%s: %s\n", Red(cresult.Name), Red(cresult.Result))
		}
	}
}

func buildjson(w io.Writer, uresults []*tester.TestUnitResult) {
	builder := simplejson.NewArrayBuilder(w)
	for _, uresult := range uresults {
		builder.AddObject(func(obj simplejson.Object) {
			obj.SetString("name", uresult.Name)
			obj.SetString("test_method", uresult.TestMethod)
			obj.SetArray("test_targets", func(arr simplejson.Array) {
				for _, tresult := range uresult.TestTargets {
					arr.AddObject(func(obj simplejson.Object) {
						obj.SetString("name", tresult.Name)
						obj.SetString("language", tresult.Language)
						obj.SetArray("test_cases", func(arr simplejson.Array) {
							for _, cresult := range tresult.TestCases {
								arr.AddObject(func(obj simplejson.Object) {
									obj.SetString("name", cresult.Name)
									obj.SetString("time", cresult.Time)
									obj.SetString("expect", cresult.Expect)
									obj.SetString("result", cresult.Result)
									// write output from file
									if fp, err := os.Open(cresult.Output); err == nil {
										obj.SetStringFromReader("output", bufio.NewReader(fp))
									}
									// write other info from files
									buf := bytes.NewBuffer([]byte{})
									for _, file := range cresult.Others {
										if fp, err := os.Open(file); err == nil {
											bufio.NewReader(fp).WriteTo(buf)
										}
									}
									obj.SetStringFromReader("others", buf)
									// errors (in this application)
									obj.SetString("errors", cresult.Errors)
								})
							}
						})
					})
				}
			})
		})
	}
	builder.Flush()
}

func printjson(uresults []*tester.TestUnitResult) {
	if prettyprint {
		writer := bytes.NewBuffer([]byte{})
		buildjson(writer, uresults)

		o := make([]struct {
			Name        string `json:"name"`
			TestMethod  string `json:"test_method"`
			TestTargets []struct {
				Name      string `json:"name"`
				Language  string `json:"language"`
				TestCases []struct {
					Name   string `json:"name"`
					Time   string `json:"time"`
					Expect string `json:"expect"`
					Result string `json:"result"`
					Output string `json:"output"`
					Others string `json:"others"`
					Errors string `json:"errors"`
				} `json:"test_cases"`
			} `json:"test_targets"`
		}, 0)

		dec := json.NewDecoder(writer)
		dec.Decode(&o)

		stdoutbuf := bufio.NewWriter(os.Stdout)

		enc := json.NewEncoder(stdoutbuf)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)

		if err := enc.Encode(o); err != nil {
			panic(err)
		}
	} else {
		buildjson(os.Stdout, uresults)
	}

}
