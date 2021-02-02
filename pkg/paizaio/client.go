package paizaio

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Host   string
	APIKey string
}

func NewClient(config Config) *Client {
	return &Client{
		Host:   config.Host,
		APIKey: config.APIKey,
	}
}

func (c *Client) api(method, endpoint string, body Request, target Response) error {
	url := fmt.Sprintf("%s/%s?api_key=%s", c.Host, endpoint, c.APIKey)
	req, err := http.NewRequest(method, url, body.Reader())
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
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

	if err := target.Write(bufio.NewReader(res.Body)); err != nil {
		return err
	}

	return nil

}

func (c *Client) RunnersCreate(req *RunnersCreateRequest) (*RunnersCreateResponse, error) {
	res := &RunnersCreateResponse{}
	if err := c.api("POST", "runners/create", req, res); err != nil {
		return res, err
	}

	if res.ID == "" {
		return res, errors.New(res.Error)
	}

	return res, nil
}

func (c *Client) RunnersGetStatus(req *RunnersGetStatusRequest) (*RunnersGetStatusResponse, error) {
	res := &RunnersGetStatusResponse{}
	if err := c.api("GET", "runners/get_status", req, res); err != nil {
		return res, err
	}

	return res, nil
}

func (c *Client) RunnersGetDetails(req *RunnersGetDetailsRequest) (*RunnersGetDetailsResponse, error) {
	res := &RunnersGetDetailsResponse{}
	if err := c.api("GET", "/runners/get_details", req, res); err != nil {
		return res, err
	}

	return res, nil
}
