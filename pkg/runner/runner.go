package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cwd-k2/titania.go/pkg/simplejson"
)

type Runner struct {
	Host   string
	APIKey string
}

func NewRunner(config Config) *Runner {
	return &Runner{
		Host:   config.Host,
		APIKey: config.APIKey,
	}
}

func (r *Runner) api(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s?api_key=%s", r.Host, endpoint, r.APIKey)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func (r *Runner) create(ospec *OrderSpec) (string, error) {
	pr, pw := io.Pipe()

	go func() {
		obj := simplejson.NewObjectBuilder(pw)
		obj.SetString("language", ospec.Language)
		obj.SetBool("longpoll", true)
		obj.SetInt("longpoll_timeout", 16)
		obj.SetStringFromReader("source_code", ospec.SourceCode)
		obj.SetStringFromReader("input", ospec.Input)
		obj.Flush()
		pw.Close()
	}()

	res, err := r.api("POST", "/runners/create", pr)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		byteArray, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		if res.StatusCode >= 500 {
			return "", &ClientError{res.StatusCode, string(byteArray)}
		} else {
			return "", &ServerError{res.StatusCode, string(byteArray)}
		}
	}

	resp := &RunnersCreateResponse{}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return "", nil
	}

	if resp.Error != "" {
		return "", &RunnerError{resp.Error}
	}

	return resp.ID, nil
}

// TODO: refactoring
func (r *Runner) Run(ospec *OrderSpec) (*ResultSpec, error) {
	id, err := r.create(ospec)
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	go func() {
		obj := simplejson.NewObjectBuilder(pw)
		obj.SetString("id", id)
		obj.Flush()
		pw.Close()
	}()

	res, err := r.api("GET", "/runners/get_details", pr)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		byteArray, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if res.StatusCode >= 500 {
			return nil, &ClientError{res.StatusCode, string(byteArray)}
		} else {
			return nil, &ServerError{res.StatusCode, string(byteArray)}
		}
	}

	rspec := &ResultSpec{}

	decoder := json.NewDecoder(res.Body)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if _, ok := token.(json.Delim); ok {
			continue
		}

		key, ok := token.(string)
		if !ok {
			return nil, errors.New("JSON parse error.")
		}

		nextToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		if num, ok := nextToken.(float64); ok {
			switch key {
			case "exit_code":
				rspec.ExitCode = int(num)
			case "build_exit_code":
				rspec.BuildExitCode = int(num)
			case "memory":
				rspec.Memory = int(num)
			case "build_memory":
				rspec.BuildMemory = int(num)
			}
		} else if str, ok := nextToken.(string); ok {
			switch key {
			case "error":
				return nil, &RunnerError{str}
			case "result":
				rspec.Result = str
			case "build_result":
				rspec.BuildResult = str
			case "time":
				rspec.Time = str
			case "build_time":
				rspec.BuildTime = str
			case "stdout":
				if _, err := strings.NewReader(str).WriteTo(ospec.Stdout); err != nil {
					return nil, err
				}
			case "stderr":
				if _, err := strings.NewReader(str).WriteTo(ospec.Stderr); err != nil {
					return nil, err
				}
			case "build_stdout":
				if _, err := strings.NewReader(str).WriteTo(ospec.BuildStdout); err != nil {
					return nil, err
				}
			case "build_stderr":
				if _, err := strings.NewReader(str).WriteTo(ospec.BuildStderr); err != nil {
					return nil, err
				}
			}
		}
	}

	return rspec, nil
}
