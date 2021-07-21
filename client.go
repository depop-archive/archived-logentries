package logentries

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	ApiKey   string
	ProxyUrl string
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

const logentriesApi = "https://rest.logentries.com/"
const logentriesLogsResource = "management/logs/"
const logentriesLogsetsResource = "management/logsets/"

func New(apikey string, proxyurlOpt ...string) *Client {
	proxyUrl := ""
	if len(proxyurlOpt) > 0 {
		proxyUrl = proxyurlOpt[0]
	}
	return &Client{
		LogSet: &LogSetClient{
			Request{
				ApiKey:   apikey,
				ProxyUrl: proxyUrl,
			},
		},
		LogSets: &LogSetsClient{
			Request{
				ApiKey:   apikey,
				ProxyUrl: proxyUrl,
			},
		},
		Log: &LogClient{
			Request{
				ApiKey:   apikey,
				ProxyUrl: proxyUrl,
			},
		},
	}
}

func (r *Request) getUrl(resource string) string {
	if r.ProxyUrl != "" {
		return r.ProxyUrl + resource
	}
	return logentriesApi + resource
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
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
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
