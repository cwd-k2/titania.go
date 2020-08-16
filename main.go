package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cwd-k2/titania.go/pretty"
	"github.com/cwd-k2/titania.go/tester"
)

const VERSION = "0.0.0-alpha"

func main() {

	// オプション解析
	// テストする言語を指定する
	var flagLanguages string
	flag.StringVar(&flagLanguages, "lang", "", "languages to test (ex. --lang=ruby,python3,java)")

	flag.Parse()
	lang := strings.Split(flagLanguages, ",")
	args := flag.Args()

	var directories []string
	if len(args) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		entries, err := ioutil.ReadDir(pwd)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				directories = append(directories, entry.Name())
			}
		}

	} else {
		directories = args
	}

	i := 0
	details := make(map[string]interface{})

	for _, dirname := range directories {
		testRoom := tester.NewTestRoom(dirname, lang)
		// 実行するテストがない
		if testRoom == nil {
			continue
		}

		i++
		fmt.Fprintf(os.Stderr, "%s\n", pretty.Bold(pretty.Green(dirname)))
		results := testRoom.Exec()

		defer tester.WrapUp(results)
		details[dirname] = results
	}

	if i == 0 {
		// 何もテストが実行されなかった場合
		println("Uh, OK, there's no test.")
	} else {
		output, err := json.MarshalIndent(details, "", "  ")
		if err != nil {
			println(err)
		}
		// 実行結果を JSON 形式で出力
		fmt.Println(string(output))
	}

}
