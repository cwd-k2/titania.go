package runner

import "io"

type OrderSpec struct {
	Language   string
	SourceCode io.Reader
	Input      io.Reader

	BuildStdout io.Writer
	BuildStderr io.Writer
	Stdout      io.Writer
	Stderr      io.Writer
}

type ResultSpec struct {
	BuildExitCode int
	ExitCode      int
	BuildTime     string
	BuildMemory   int
	BuildResult   string
	Time          string
	Memory        int
	Result        string
}
