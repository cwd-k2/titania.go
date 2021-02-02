package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"

	. "github.com/cwd-k2/titania.go/internal/pkg/pretty"
	"github.com/cwd-k2/titania.go/pkg/tester"
)

func exec(directories []string) {
	var (
		// buffer, buffer, buffering!
		buffer   = bytes.NewBuffer([]byte{})
		buferr   = bytes.NewBuffer([]byte{})
		bufout   = bytes.NewBuffer([]byte{})
		stdout   = bufio.NewWriter(os.Stdout)
		stderr   = bufio.NewWriter(os.Stderr)
		bufenc   = json.NewEncoder(buffer)
		executed = 0
	)
	defer stdout.Flush()
	defer stderr.Flush()

	// don't escape for HTML emmbedding
	bufenc.SetEscapeHTML(false)

	// open square bracket
	buffer.WriteByte('[')

	for _, dirname := range directories {
		// 設定
		tconf := tester.NewConfig(dirname)
		if tconf == nil {
			continue
		}
		// TestUnit
		tunit := tester.NewTestUnit(dirname, tconf)
		if tunit == nil {
			continue
		}

		outcome := tunit.Exec()

		// need commas.
		if executed != 0 {
			buffer.WriteByte(',')
		}

		printoutcome(buferr, outcome)

		// store as json bytes instead of structs themselves
		if err := bufenc.Encode(outcome); err != nil {
			panic(err)
		}

		executed++
	}

	// closing square bracket
	buffer.WriteByte(']')

	if executed == 0 {
		// 何もテストが実行されなかった場合
		println("There's no test in (sub)directory[ies].")
		os.Exit(1)
	} else {
		// all done.
		Printf("\n%s\n", Bold("ALL DONE"))

		// output buffered info
		if _, err := buferr.WriteTo(stderr); err != nil {
			panic(err)
		}

		// reshape buffered json
		if err := json.Indent(bufout, buffer.Bytes(), "", "  "); err != nil {
			panic(err)
		}

		// now output buffered json data
		if _, err := bufout.WriteTo(stdout); err != nil {
			panic(err)
		}
	}
}
