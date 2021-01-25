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
	http   http.Client
}

func NewClient(config Config) *Client {
	return &Client{config.Host, config.APIKey, http.Client{}}
}

func (c *Client) Api(method, endpoint string, params map[string]string, target interface{}) error {

	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		return nil
	}

	req, err := http.NewRequest(method, c.Host+endpoint, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		byteArray, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode >= 500 {
			return &ClientError{res.StatusCode, string(byteArray)}
		} else {
			return &ServerError{res.StatusCode, string(byteArray)}
		}
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return err
	}

	return nil

}

func (c *Client) RunnersCreate(language string, sourceCode, input string) (RunnersCreateResponse, error) {
	args := map[string]string{
		"api_key":          c.APIKey,
		"language":         language,
		"source_code":      sourceCode,
		"input":            input,
		"longpoll":         "true",
		"longpoll_timeout": "30",
	}

	var res RunnersCreateResponse

	if err := c.Api("POST", "/runners/create", args, &res); err != nil {
		return res, err
	}

	if res.ID == "" {
		return res, errors.New(res.Error)
	}

	return res, nil
}

func (c *Client) RunnersGetDetails(id string) (RunnersGetDetailsResponse, error) {
	args := map[string]string{
		"api_key": c.APIKey,
		"id":      id,
	}

	var res RunnersGetDetailsResponse

	if err := c.Api("GET", "/runners/get_details", args, &res); err != nil {
		return res, err
	}

	return res, nil
}
