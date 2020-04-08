package logentries

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	ApiKey string
}
type LogSetClient struct {
	Request
}
type LogSetsClient struct {
	Request
}
type LogClient struct {
	Request
}

type Client struct {
	Log     *LogClient
	LogSet  *LogSetClient
	LogSets *LogSetsClient
}

func New(apikey string) *Client {
	return &Client{
		LogSet: &LogSetClient{
			Request{
				ApiKey: apikey,
			},
		},
		LogSets: &LogSetsClient{
			Request{
				ApiKey: apikey,
			},
		},
		Log: &LogClient{
			Request{
				ApiKey: apikey,
			},
		},
	}
}

func (r *Request) getLogentries(url string, expected int) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return r.performRequest(req, expected)
}

func (r *Request) postLogentries(url string, payload []byte, expected int) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	return r.performRequest(req, expected)
}

func (r *Request) putLogentries(url string, payload []byte, expected int) ([]byte, error) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	return r.performRequest(req, expected)
}
func (r *Request) deleteLogentries(url string, expected int) (bool, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return false, err
	}

	_, err = r.performRequest(req, expected)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *Request) performRequest(req *http.Request, expected int) ([]byte, error) {
	client := &http.Client{Timeout: time.Duration(5) * time.Second}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", r.ApiKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if res.Body != nil {
		body, err := ioutil.ReadAll(res.Body)
		if res.StatusCode != expected {
			return nil, fmt.Errorf("unexpected logentries response code: %v, payload: %s", res.StatusCode, body)
		}
		return body, err
	}

	if res.StatusCode != expected {
		return nil, fmt.Errorf("unexpected logentries response code: %v", res.StatusCode)
	}

	return nil, nil
}
