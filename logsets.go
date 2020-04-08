package logentries

import (
	"encoding/json"
	"net/http"
)

type LogsetsResponse struct {
	Logsets []LogSet `json:"logsets"`
}

type LogSetsReadResponse struct {
	LogSets []LogSet `json:"logsets"`
}

type LogSetsReadRequest struct{}

func (l *LogSetsClient) Read(readRequest *LogSetsReadRequest) (LogSetsReadResponse, error) {
	url := "https://rest.logentries.com/management/logsets"

	resp, err := l.getLogentries(url, http.StatusOK)
	if err != nil {
		return LogSetsReadResponse{}, err
	}

	var logsets LogSetsReadResponse
	err = json.Unmarshal(resp, &logsets)
	if err != nil {
		return LogSetsReadResponse{}, err
	}

	return logsets, nil
}
