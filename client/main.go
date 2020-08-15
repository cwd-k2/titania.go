package client

import (
	"encoding/json"
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
		return nil, err
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
		return nil, err
	}

	return runnersGetDetailsResponse, nil
}
