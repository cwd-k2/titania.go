package runner

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
