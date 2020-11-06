package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type RunnersCreateResponse struct {
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
	BuildExitCode uint   `json:"build_exit_code"`
	BuildTime     string `json:"bulid_time"`
	BuildMemory   uint   `json:"build_memory"`
	BuildResult   string `json:"build_result"`
	STDOUT        string `json:"stdout"`
	STDERR        string `json:"stderr"`
	ExitCode      uint   `json:"exit_code"`
	Time          string `json:"time"`
	Memory        uint   `json:"memory"`
	Connections   uint   `json:"connections"`
	Result        string `json:"result"`
}

type Client struct {
	Host   string
	APIKey string
}

func NewClient(config Config) *Client {
	return &Client{config.Host, config.APIKey}
}

func (c *Client) Api(method, endpoint string, params map[string]string, target interface{}) *TitaniaClientError {

	body, err := json.Marshal(params)
	if err != nil {
		return &TitaniaClientError{-1, err}
	}

	request, err := http.NewRequest(method, c.Host+endpoint, bytes.NewReader(body))
	if err != nil {
		return &TitaniaClientError{-1, err}
	}

	request.Header.Set("Content-Type", "application/json")

	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		return &TitaniaClientError{-1, err}
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		byteArray, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return &TitaniaClientError{-1, err}
		}

		return &TitaniaClientError{response.StatusCode, errors.New(string(byteArray))}
	}

	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		return &TitaniaClientError{-1, err}
	}

	return nil

}

func (c *Client) RunnersCreate(language string, sourceCode, input *bytes.Buffer) (*RunnersCreateResponse, *TitaniaClientError) {
	args := map[string]string{
		"api_key":          c.APIKey,
		"language":         language,
		"source_code":      sourceCode.String(),
		"input":            input.String(),
		"longpoll":         "true",
		"longpoll_timeout": "30",
	}

	runnersCreateResponse := new(RunnersCreateResponse)

	if err := c.Api("POST", "/runners/create", args, runnersCreateResponse); err != nil {
		return nil, err
	}

	if runnersCreateResponse.ID == "" {
		return nil, &TitaniaClientError{-1, errors.New(runnersCreateResponse.Error)}
	}

	return runnersCreateResponse, nil
}

func (c *Client) RunnersGetDetails(id string) (*RunnersGetDetailsResponse, *TitaniaClientError) {
	args := map[string]string{
		"api_key": c.APIKey,
		"id":      id,
	}

	runnersGetDetailsResponse := new(RunnersGetDetailsResponse)

	if err := c.Api("GET", "/runners/get_details", args, runnersGetDetailsResponse); err != nil {
		return nil, err
	}

	return runnersGetDetailsResponse, nil
}
