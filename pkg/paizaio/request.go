package paizaio

import (
	"bytes"
	"encoding/json"
	"io"
)

type Request interface {
	Reader() io.Reader
}

type RunnersCreateRequest struct {
	Language        string `json:"language"`
	SourceCode      string `json:"source_code"`
	Input           string `json:"input"`
	Longpoll        bool   `json:"longpoll"`
	LongpollTimeout int    `json:"longpoll_timeout"`
}

type RunnersGetStatusRequest struct {
	ID string `json:"id"`
}

type RunnersGetDetailsRequest struct {
	ID string `json:"id"`
}

func (r *RunnersCreateRequest) Reader() io.Reader {
	p := bytes.NewBuffer([]byte{})
	json.NewEncoder(p).Encode(r)
	return p
}

func (r *RunnersGetStatusRequest) Reader() io.Reader {
	p := bytes.NewBuffer([]byte{})
	json.NewEncoder(p).Encode(r)
	return p
}

func (r *RunnersGetDetailsRequest) Reader() io.Reader {
	p := bytes.NewBuffer([]byte{})
	json.NewEncoder(p).Encode(r)
	return p
}
