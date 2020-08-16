package client

import (
	"encoding/json"
	"errors"
	"fmt"
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
	APIKey string
	Host   string
}

type ClientError struct {
	Err error
}

type ServerError struct {
	Err error
}

func (e ClientError) Error() string {
	return e.Err.Error()
}

func (e ServerError) Error() string {
	return e.Err.Error()
}

func (c *Client) api(
	method string,
	endpoint string,
	params map[string]string) ([]byte, error) {

	httpClient := new(http.Client)

	request, err := http.NewRequest(method, c.Host+endpoint, nil)
	if err != nil {
		return nil, err
	}

	query := request.URL.Query()
	query.Add("api_key", c.APIKey)
	for k, v := range params {
		query.Add(k, v)
	}

	request.URL.RawQuery = query.Encode()

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	byteArray, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 500 {
		return nil, ServerError{errors.New(string(byteArray))}
	}

	if response.StatusCode >= 400 {
		return nil, ClientError{errors.New(string(byteArray))}
	}

	return byteArray, nil

}

func (c *Client) RunnersCreate(
	sourceCode string,
	language string,
	input string) (*RunnersCreateResponse, error) {

	runnersCreateResponse := new(RunnersCreateResponse)

	args := make(map[string]string)
	args["source_code"] = sourceCode
	args["language"] = language
	args["input"] = input
	args["longpoll"] = "true"
	args["longpoll_timeout"] = "100"

	byteArray, err := c.api("POST", "/runners/create", args)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteArray, runnersCreateResponse); err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", err.Error(), string(byteArray)))
	}

	return runnersCreateResponse, nil
}

func (c *Client) RunnersGetDetails(
	id string) (*RunnersGetDetailsResponse, error) {

	runnersGetDetailsResponse := new(RunnersGetDetailsResponse)

	args := make(map[string]string)
	args["id"] = id

	byteArray, err := c.api("GET", "/runners/get_details", args)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteArray, runnersGetDetailsResponse); err != nil {
		return nil, errors.New(fmt.Sprintf("%s\n%s", err.Error(), string(byteArray)))
	}

	return runnersGetDetailsResponse, nil
}
