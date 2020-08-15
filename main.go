package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cwd-k2/titania.go/pretty"
	"github.com/cwd-k2/titania.go/tester"
)

const VERSION = "0.0.0-alpha"

func main() {
	args := os.Args[1:]
	var directories []string
	// TODO: オプション解析したいね
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

	testRooms := tester.MakeTestRooms(directories)
	for roomName, testRoom := range testRooms {
		fmt.Printf("%s\n", pretty.Bold(pretty.Green(roomName)))
		testRoom.Exec()
	}

}
