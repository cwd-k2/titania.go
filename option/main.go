package option

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
)

// オプション解析
func Parse() ([]string, []string) {
	// テストする言語を指定する
	var flagLanguages string
	flag.StringVar(
		&flagLanguages,
		"lang", "", "languages to test (ex. --lang=ruby,python3,java)")

	flag.Parse()
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

	return directories, strings.Split(flagLanguages, ",")
}
