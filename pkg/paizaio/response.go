package paizaio

import (
	"encoding/json"
	"io"
)

type Response interface {
	Write(io.Reader) error
}

type RunnersCreateResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

type RunnersGetStatusResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

type RunnersGetDetailsResponse struct {
	ID            string `json:"id"`
	Language      string `json:"language"`
	Note          string `json:"note"`
	Status        string `json:"status"`
	BuildSTDOUT   string `json:"build_stdout"`
	BuildSTDERR   string `json:"build_stderr"`
	BuildExitCode int    `json:"build_exit_code"`
	BuildTime     string `json:"build_time"`
	BuildMemory   int    `json:"build_memory"`
	BuildResult   string `json:"build_result"`
	STDOUT        string `json:"stdout"`
	STDERR        string `json:"stderr"`
	ExitCode      int    `json:"exit_code"`
	Time          string `json:"time"`
	Memory        int    `json:"memory"`
	Connections   int    `json:"connections"`
	Result        string `json:"result"`
	Error         string `json:"error"`
}

func (r *RunnersCreateResponse) Write(w io.Reader) error {
	return json.NewDecoder(w).Decode(r)
}

func (r *RunnersGetStatusResponse) Write(w io.Reader) error {
	return json.NewDecoder(w).Decode(r)
}

func (r *RunnersGetDetailsResponse) Write(w io.Reader) error {
	return json.NewDecoder(w).Decode(r)
}
